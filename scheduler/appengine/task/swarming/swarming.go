// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package swarming implements tasks that run Swarming jobs.
package swarming

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/api/pubsub/v1"

	"github.com/luci/gae/service/info"
	"github.com/luci/luci-go/common/api/swarming/swarming/v1"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/scheduler/appengine/messages"
	"github.com/luci/luci-go/scheduler/appengine/task"
	"github.com/luci/luci-go/scheduler/appengine/task/utils"
)

const (
	statusCheckTimerName     = "check-swarming-task-status"
	statusCheckTimerInterval = time.Minute
)

// TaskManager implements task.Manager interface for tasks defined with
// SwarmingTask proto message.
type TaskManager struct {
}

// Name is part of Manager interface.
func (m TaskManager) Name() string {
	return "swarming"
}

// ProtoMessageType is part of Manager interface.
func (m TaskManager) ProtoMessageType() proto.Message {
	return (*messages.SwarmingTask)(nil)
}

// ValidateProtoMessage is part of Manager interface.
func (m TaskManager) ValidateProtoMessage(msg proto.Message) error {
	cfg, ok := msg.(*messages.SwarmingTask)
	if !ok {
		return fmt.Errorf("wrong type %T, expecting *messages.SwarmingTask", msg)
	}
	if cfg == nil {
		return fmt.Errorf("expecting a non-empty SwarmingTask")
	}

	// Validate 'server' field.
	if cfg.Server == "" {
		return fmt.Errorf("field 'server' is required")
	}
	u, err := url.Parse(cfg.Server)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %s", cfg.Server, err)
	}
	if !u.IsAbs() {
		return fmt.Errorf("not an absolute url: %q", cfg.Server)
	}
	if u.Path != "" {
		return fmt.Errorf("not a host root url: %q", cfg.Server)
	}

	// Validate environ, dimensions, tags.
	if err = utils.ValidateKVList("environment variable", cfg.Env, '='); err != nil {
		return err
	}
	if err = utils.ValidateKVList("dimension", cfg.Dimensions, ':'); err != nil {
		return err
	}
	if err = utils.ValidateKVList("tag", cfg.Tags, ':'); err != nil {
		return err
	}

	// Default tags can not be overridden.
	defTags := defaultTags(nil, nil)
	for _, kv := range utils.UnpackKVList(cfg.Tags, ':') {
		if _, ok := defTags[kv.Key]; ok {
			return fmt.Errorf("tag %q is reserved", kv.Key)
		}
	}

	// Validate priority.
	priority := cfg.Priority
	if priority != 0 && (priority < 0 || priority > 255) {
		return fmt.Errorf("bad priority, must be (0, 255]: %d", priority)
	}

	// Can't have both 'command' and 'isolated_ref'.
	hasCommand := len(cfg.Command) != 0
	hasIsolatedRef := cfg.IsolatedRef != nil
	switch {
	case !hasCommand && !hasIsolatedRef:
		return fmt.Errorf("one of 'command' or 'isolated_ref' is required")
	case hasCommand && hasIsolatedRef:
		return fmt.Errorf("only one of 'command' or 'isolated_ref' must be specified, not both")
	}

	return nil
}

func kvListToStringPairs(list []string, sep rune) (out []*swarming.SwarmingRpcsStringPair) {
	for _, pair := range utils.UnpackKVList(list, sep) {
		out = append(out, &swarming.SwarmingRpcsStringPair{
			Key:   pair.Key,
			Value: pair.Value,
		})
	}
	return out
}

// defaultTags returns map with default set of tags.
//
// If context is nil, only keys are set.
func defaultTags(c context.Context, ctl task.Controller) map[string]string {
	if c != nil {
		return map[string]string{
			"scheduler_invocation_id": fmt.Sprintf("%d", ctl.InvocationID()),
			"scheduler_job_id":        ctl.JobID(),
			"user_agent":              info.AppID(c),
		}
	}
	return map[string]string{
		"scheduler_invocation_id": "",
		"scheduler_job_id":        "",
		"user_agent":              "",
	}
}

// defaultExpirationTimeout derives Swarming queuing timeout: max time a task
// is kept in the queue (not being picked up by bots), before it is marked as
// failed.
func defaultExpirationTimeout(ctl task.Controller) time.Duration {
	// TODO(vadimsh): Do something smarter, e.g. look at next expected invocation
	// time.
	return 30 * time.Minute
}

// defaultExecutionTimeout derives hard deadline for a task if it wasn't
// explicitly specified in the config.
func defaultExecutionTimeout(ctl task.Controller) time.Duration {
	// TODO(vadimsh): Do something smarter, e.g. look at next expected invocation
	// time.
	return time.Hour
}

// taskData is saved in Invocation.TaskData field.
type taskData struct {
	SwarmingTaskID string `json:"swarming_task_id"`
}

// LaunchTask is part of Manager interface.
func (m TaskManager) LaunchTask(c context.Context, ctl task.Controller) error {
	// At this point config is already validated by ValidateProtoMessage.
	cfg := ctl.Task().(*messages.SwarmingTask)

	// Default set of tags.
	tags := utils.KVListFromMap(defaultTags(c, ctl)).Pack(':')
	tags = append(tags, cfg.Tags...)

	// How long to keep a task in swarming queue (not running) before marking it
	// as expired.
	expirationSecs := int64(defaultExpirationTimeout(ctl) / time.Second)

	// The hard deadline: how long task can run once it has started.
	executionTimeoutSecs := int64(cfg.ExecutionTimeoutSecs)
	if executionTimeoutSecs == 0 {
		executionTimeoutSecs = int64(defaultExecutionTimeout(ctl) / time.Second)
	}

	// How long the task is allowed to spend in cleanup.
	gracePeriodSecs := cfg.GracePeriodSecs
	if gracePeriodSecs == 0 {
		gracePeriodSecs = 30
	}

	// The task priority (or niceness, lower is more aggressive).
	priority := cfg.Priority
	if priority == 0 {
		priority = 200
	}

	// Make sure Swarming can publish PubSub messages, grab token that would
	// identify this invocation when receiving PubSub notifications.
	ctl.DebugLog("Preparing PubSub topic for %q", cfg.Server)
	topic, authToken, err := ctl.PrepareTopic(cfg.Server)
	if err != nil {
		ctl.DebugLog("Failed to prepare PubSub topic - %s", err)
		return err
	}
	ctl.DebugLog("PubSub topic is %q", topic)

	// Prepare the request.
	request := swarming.SwarmingRpcsNewTaskRequest{
		Name:            fmt.Sprintf("scheduler:%s/%d", ctl.JobID(), ctl.InvocationID()),
		ExpirationSecs:  expirationSecs,
		Priority:        int64(priority),
		PubsubAuthToken: "...", // set a bit later, after printing this struct
		PubsubTopic:     topic,
		Tags:            tags,
		Properties: &swarming.SwarmingRpcsTaskProperties{
			Dimensions:           kvListToStringPairs(cfg.Dimensions, ':'),
			Env:                  kvListToStringPairs(cfg.Env, '='),
			ExecutionTimeoutSecs: executionTimeoutSecs,
			ExtraArgs:            cfg.ExtraArgs,
			GracePeriodSecs:      int64(gracePeriodSecs),
			Idempotent:           false,
			IoTimeoutSecs:        int64(cfg.IoTimeoutSecs),
		},
	}

	// Only one of InputsRef or Command must be set.
	if cfg.IsolatedRef != nil {
		request.Properties.InputsRef = &swarming.SwarmingRpcsFilesRef{
			Isolated:       cfg.IsolatedRef.Isolated,
			Isolatedserver: cfg.IsolatedRef.IsolatedServer,
			Namespace:      cfg.IsolatedRef.Namespace,
		}
	} else {
		request.Properties.Command = cfg.Command
	}

	// Serialize for debug log without auth token.
	blob, err := json.MarshalIndent(&request, "", "  ")
	if err != nil {
		return err
	}
	ctl.DebugLog("Swarming task request:\n%s", blob)
	request.PubsubAuthToken = authToken // can put the token now

	// The next call may take a while. Dump the current log to the datastore.
	// Ignore errors here, it is best effort attempt to update the log.
	ctl.Save()

	// Trigger the task.
	service, err := m.createSwarmingService(c, ctl)
	if err != nil {
		return err
	}
	resp, err := service.Tasks.New(&request).Do()
	if err != nil {
		ctl.DebugLog("Failed to trigger the task - %s", err)
		return utils.WrapAPIError(err)
	}

	// Dump response in full to the debug log. It doesn't contain any secrets.
	blob, err = json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	ctl.DebugLog("Swarming response:\n%s", blob)

	// Save TaskId in invocation, will be used later when handling PubSub
	// notifications
	ctl.State().TaskData, err = json.Marshal(&taskData{
		SwarmingTaskID: resp.TaskId,
	})
	if err != nil {
		return err
	}

	// Successfully launched.
	ctl.State().Status = task.StatusRunning
	ctl.State().ViewURL = fmt.Sprintf("%s/user/task/%s", cfg.Server, resp.TaskId)
	ctl.DebugLog("Task URL: %s", ctl.State().ViewURL)

	// Maybe the task was already finished? Can only happen when 'idempotent' is
	// set to true (which we don't do currently), but handle this case here for
	// completeness anyway.
	if resp.TaskResult != nil {
		ctl.DebugLog("Task request was deduplicated")
		m.handleTaskResult(c, ctl, resp.TaskResult)
	}

	// This will schedule status check if the task is actually running.
	m.checkTaskStatusLater(c, ctl)
	return nil
}

// AbortTask is part of Manager interface.
func (m TaskManager) AbortTask(c context.Context, ctl task.Controller) error {
	// TODO(vadimsh): Send the abort signal to Swarming.
	return nil
}

// HandleNotification is part of Manager interface.
func (m TaskManager) HandleNotification(c context.Context, ctl task.Controller, msg *pubsub.PubsubMessage) error {
	ctl.DebugLog("Received PubSub notification, asking swarming for the task status")
	return m.checkTaskStatus(c, ctl)
}

// HandleTimer is part of Manager interface.
func (m TaskManager) HandleTimer(c context.Context, ctl task.Controller, name string, payload []byte) error {
	if name == statusCheckTimerName {
		ctl.DebugLog("Timer tick, asking swarming for the task status")
		if err := m.checkTaskStatus(c, ctl); err != nil {
			// This is either a fatal or transient error. If it is fatal, no need to
			// schedule the timer anymore. If it is transient, HandleTimer call itself
			// will be retried and the timer when be rescheduled then.
			return err
		}
		m.checkTaskStatusLater(c, ctl) // reschedule this check
	}
	return nil
}

// createSwarmingService makes a configured Swarming API client.
func (m TaskManager) createSwarmingService(c context.Context, ctl task.Controller) (*swarming.Service, error) {
	client, err := ctl.GetClient(time.Minute)
	if err != nil {
		return nil, err
	}
	service, err := swarming.New(client)
	if err != nil {
		return nil, err
	}
	cfg := ctl.Task().(*messages.SwarmingTask)
	service.BasePath = cfg.Server + "/_ah/api/swarming/v1/"
	return service, nil
}

// checkTaskStatusLater schedules a delayed call to checkTaskStatus if the
// invocation is still running.
//
// This is a fallback mechanism in case PubSub notifications are delayed or
// lost for some reason.
func (m TaskManager) checkTaskStatusLater(c context.Context, ctl task.Controller) {
	// TODO(vadimsh): Make the check interval configurable?
	if !ctl.State().Status.Final() {
		ctl.AddTimer(statusCheckTimerInterval, statusCheckTimerName, nil)
	}
}

// checkTaskStatus fetches current task status from Swarming and update the
// invocation status based on it.
//
// Called on PubSub notifications and also periodically by timer.
func (m TaskManager) checkTaskStatus(c context.Context, ctl task.Controller) error {
	switch status := ctl.State().Status; {
	// This can happen if Swarming manages to send PubSub message before
	// LaunchTask finishes. Do not touch State or DebugLog to avoid collision with
	// still running LaunchTask when saving the invocation, it will only make the
	// matters worse.
	case status == task.StatusStarting:
		return errors.WrapTransient(errors.New("invocation is still starting, try again later"))
	case status != task.StatusRunning:
		return fmt.Errorf("unexpected invocation status %q, expecting %q", status, task.StatusRunning)
	}

	// Grab task ID from the blob generated in LaunchTask.
	taskData := taskData{}
	if err := json.Unmarshal(ctl.State().TaskData, &taskData); err != nil {
		ctl.State().Status = task.StatusFailed
		return fmt.Errorf("could not parse TaskData - %s", err)
	}

	// Fetch task result from Swarming.
	service, err := m.createSwarmingService(c, ctl)
	if err != nil {
		return err
	}
	resp, err := service.Task.Result(taskData.SwarmingTaskID).Do()
	if err != nil {
		ctl.DebugLog("Failed to fetch task results - %s", err)
		err = utils.WrapAPIError(err)
		if !errors.IsTransient(err) {
			ctl.State().Status = task.StatusFailed
		}
		return err
	}

	blob, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	m.handleTaskResult(c, ctl, resp)

	// Dump final status in full to the debug log. It doesn't contain any secrets.
	if ctl.State().Status.Final() {
		ctl.DebugLog("Swarming task status:\n%s", blob)
	}

	return nil
}

// handleTaskResult processes swarming task results message updating the state
// of the invocation.
func (m TaskManager) handleTaskResult(c context.Context, ctl task.Controller, r *swarming.SwarmingRpcsTaskResult) {
	ctl.DebugLog(
		"The task is in state %q (failure: %v, internalFailure: %v)",
		r.State, r.Failure, r.InternalFailure)
	switch {
	case r.State == "PENDING" || r.State == "RUNNING":
		return // do nothing
	case r.State == "COMPLETED" && !(r.Failure || r.InternalFailure):
		ctl.State().Status = task.StatusSucceeded
	default:
		ctl.State().Status = task.StatusFailed
	}
}
