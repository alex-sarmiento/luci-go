// AUTOGENERATED. DO NOT EDIT.

package testservices

import discovery "github.com/luci/luci-go/server/discovery"

func init() {
	discovery.RegisterDescriptorSetCompressed(
		[]string{
			"testservices.Greeter", "testservices.Calc",
		},
		[]byte{31, 139,
			8, 0, 0, 9, 110, 136, 0, 255, 116, 148, 223, 114, 226, 54,
			20, 135, 45, 89, 24, 115, 48, 193, 81, 33, 5, 15, 38, 138,
			123, 145, 78, 103, 74, 27, 58, 237, 69, 59, 237, 76, 66, 51,
			253, 155, 94, 148, 246, 1, 28, 16, 224, 25, 99, 83, 203, 164,
			165, 175, 176, 15, 177, 247, 251, 2, 251, 122, 123, 36, 236, 108,
			118, 103, 115, 231, 159, 207, 209, 119, 62, 73, 6, 120, 1, 224,
			111, 100, 154, 230, 255, 230, 69, 186, 156, 236, 138, 188, 204, 185,
			87, 74, 85, 42, 89, 60, 36, 11, 169, 162, 17, 120, 63, 235,
			142, 63, 229, 63, 123, 124, 207, 61, 96, 89, 188, 149, 3, 34,
			200, 167, 173, 40, 4, 168, 170, 187, 244, 192, 187, 208, 220, 74,
			165, 226, 117, 93, 190, 132, 238, 221, 62, 45, 19, 44, 214, 235,
			91, 64, 254, 51, 213, 134, 126, 60, 12, 168, 126, 68, 142, 255,
			182, 81, 237, 242, 76, 73, 93, 254, 255, 216, 57, 189, 131, 230,
			79, 133, 148, 165, 44, 248, 13, 184, 243, 248, 96, 134, 242, 96,
			242, 84, 117, 242, 212, 51, 24, 124, 176, 134, 252, 200, 154, 206,
			129, 205, 226, 116, 193, 127, 3, 183, 158, 202, 195, 119, 251, 223,
			211, 14, 198, 207, 149, 143, 178, 145, 245, 235, 235, 38, 56, 156,
			49, 235, 19, 2, 175, 8, 16, 143, 219, 204, 226, 211, 151, 68,
			204, 242, 221, 161, 72, 214, 155, 82, 76, 191, 188, 250, 70, 252,
			181, 145, 98, 182, 41, 242, 109, 178, 223, 138, 235, 125, 185, 201,
			11, 53, 17, 215, 105, 42, 76, 147, 18, 133, 212, 99, 228, 114,
			2, 226, 111, 37, 69, 190, 18, 229, 38, 81, 66, 229, 251, 98,
			33, 197, 34, 95, 74, 129, 113, 157, 63, 200, 34, 147, 75, 113,
			127, 16, 177, 184, 153, 255, 248, 185, 42, 15, 169, 20, 41, 250,
			161, 16, 174, 137, 75, 177, 136, 51, 113, 47, 65, 172, 242, 125,
			182, 20, 73, 134, 111, 165, 248, 253, 151, 217, 237, 31, 243, 91,
			177, 74, 82, 57, 1, 112, 129, 80, 110, 59, 110, 15, 38, 64,
			29, 139, 179, 150, 229, 145, 32, 50, 158, 107, 125, 232, 73, 182,
			22, 213, 198, 197, 82, 174, 146, 44, 41, 147, 60, 195, 149, 0,
			182, 99, 17, 110, 183, 220, 46, 156, 3, 115, 44, 106, 113, 187,
			77, 191, 14, 184, 152, 203, 108, 169, 208, 171, 6, 0, 120, 208,
			208, 13, 216, 222, 118, 78, 234, 132, 131, 219, 254, 168, 78, 54,
			166, 203, 43, 248, 30, 40, 158, 28, 235, 90, 167, 36, 184, 50,
			26, 197, 241, 18, 68, 245, 113, 225, 25, 100, 101, 140, 30, 40,
			166, 55, 180, 71, 187, 75, 37, 244, 103, 121, 180, 98, 122, 76,
			23, 119, 212, 6, 188, 16, 109, 229, 83, 14, 29, 104, 232, 192,
			56, 243, 105, 247, 76, 15, 213, 177, 161, 139, 110, 157, 112, 157,
			223, 234, 212, 9, 133, 124, 255, 20, 190, 67, 33, 194, 89, 207,
			58, 35, 193, 23, 149, 208, 241, 218, 159, 51, 170, 183, 173, 142,
			58, 4, 177, 61, 151, 27, 29, 162, 117, 250, 180, 111, 116, 136,
			209, 233, 211, 94, 207, 140, 36, 70, 167, 95, 233, 16, 163, 211,
			111, 249, 117, 66, 157, 254, 71, 61, 36, 82, 7, 117, 6, 86,
			64, 204, 21, 104, 250, 192, 245, 52, 221, 49, 244, 33, 253, 193,
			28, 233, 17, 48, 172, 142, 155, 152, 227, 30, 250, 231, 117, 66,
			220, 240, 179, 111, 53, 142, 81, 206, 70, 214, 185, 193, 49, 138,
			107, 70, 238, 199, 70, 150, 106, 92, 72, 79, 140, 44, 53, 178,
			33, 29, 13, 141, 16, 53, 178, 33, 109, 214, 9, 215, 133, 110,
			171, 78, 72, 15, 189, 78, 69, 193, 210, 248, 145, 66, 144, 50,
			166, 225, 73, 213, 73, 26, 186, 88, 83, 244, 102, 198, 143, 20,
			130, 148, 49, 82, 180, 163, 205, 217, 133, 254, 113, 105, 71, 27,
			187, 46, 220, 129, 161, 219, 218, 49, 170, 232, 182, 113, 140, 232,
			69, 96, 8, 182, 113, 140, 42, 186, 109, 28, 163, 138, 110, 27,
			199, 200, 235, 220, 59, 230, 127, 239, 171, 55, 1, 0, 0, 255,
			255, 101, 243, 117, 100, 14, 5, 0, 0},
	)
}
