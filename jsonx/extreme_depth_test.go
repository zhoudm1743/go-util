package jsonx

import (
	"fmt"
	"testing"
)

// 测试不同深度的嵌套能力
func TestDepthCapabilities(t *testing.T) {
	depths := []int{10, 20, 50, 100, 200}

	for _, depth := range depths {
		t.Run(fmt.Sprintf("Depth_%d", depth), func(t *testing.T) {
			j := Object()

			// 构建路径
			path := "root"
			for i := 1; i <= depth; i++ {
				path += fmt.Sprintf(".level_%d", i)
			}

			// 设置值
			expectedValue := fmt.Sprintf("depth_%d_value", depth)
			j.Set(path, expectedValue)

			// 验证获取
			actualValue := j.Get(path).String()
			if actualValue != expectedValue {
				t.Errorf("深度 %d 失败: expected '%s', got '%s'", depth, expectedValue, actualValue)
				return
			}

			// 验证根路径存在
			if !j.Has("root") {
				t.Errorf("深度 %d: 根路径应该存在", depth)
				return
			}

			// 验证中间路径（如果深度足够）
			if depth >= 5 {
				midPath := "root.level_1.level_2.level_3.level_4.level_5"
				if !j.Has(midPath) {
					t.Errorf("深度 %d: 中间路径 %s 应该存在", depth, midPath)
					return
				}
			}

			t.Logf("✅ 深度 %d 层嵌套测试通过", depth)
		})
	}
}

// 测试对象数组混合的不同深度
func TestObjectArrayMixedDepth(t *testing.T) {
	j := Object()

	// 测试路径：arr[0].obj.arr[1].obj.arr[2].value
	// 这代表：数组->对象->数组->对象->数组->值
	paths := []string{
		"level1.0.data",                            // 2层：对象->数组->对象
		"level1.0.level2.1.data",                   // 4层：对象->数组->对象->数组->对象
		"level1.0.level2.1.level3.2.data",          // 6层：对象->数组->对象->数组->对象->数组->对象
		"level1.0.level2.1.level3.2.level4.3.data", // 8层：更深的混合
	}

	for i, path := range paths {
		value := fmt.Sprintf("mixed_value_%d", i)
		j.Set(path, value)

		// 验证设置成功
		if actual := j.Get(path).String(); actual != value {
			t.Errorf("混合路径 %s 失败: expected '%s', got '%s'", path, value, actual)
		}

		t.Logf("✅ 混合路径 %s 设置成功", path)
	}

	// 打印结构查看
	jsonStr, _ := j.ToPrettyJSON()
	t.Logf("混合对象数组结构:\n%s", jsonStr)
}

// 测试大型对象数组
func TestLargeObjectArray(t *testing.T) {
	j := Object()

	// 创建 users[0-99].profiles[0-4].settings[0-2].value 的结构
	userCount := 100
	profileCount := 5
	settingCount := 3

	for u := 0; u < userCount; u++ {
		j.Set(fmt.Sprintf("users.%d.id", u), u)
		j.Set(fmt.Sprintf("users.%d.name", u), fmt.Sprintf("User_%d", u))

		for p := 0; p < profileCount; p++ {
			j.Set(fmt.Sprintf("users.%d.profiles.%d.id", u, p), p)
			j.Set(fmt.Sprintf("users.%d.profiles.%d.name", u, p), fmt.Sprintf("Profile_%d_%d", u, p))

			for s := 0; s < settingCount; s++ {
				settingPath := fmt.Sprintf("users.%d.profiles.%d.settings.%d", u, p, s)
				j.Set(settingPath+".key", fmt.Sprintf("setting_%d", s))
				j.Set(settingPath+".value", fmt.Sprintf("value_%d_%d_%d", u, p, s))
			}
		}
	}

	// 验证总数据量
	totalSettings := userCount * profileCount * settingCount
	t.Logf("创建了 %d 个用户，每个用户 %d 个配置，每个配置 %d 个设置，总计 %d 个设置对象",
		userCount, profileCount, settingCount, totalSettings)

	// 验证一些随机数据
	testCases := []struct {
		path     string
		expected string
	}{
		{"users.0.name", "User_0"},
		{"users.50.profiles.2.name", "Profile_50_2"},
		{"users.99.profiles.4.settings.2.value", "value_99_4_2"},
	}

	for _, test := range testCases {
		if actual := j.Get(test.path).String(); actual != test.expected {
			t.Errorf("路径 %s 失败: expected '%s', got '%s'", test.path, test.expected, actual)
		}
	}

	// 验证数组长度
	if j.Get("users").Length() != userCount {
		t.Errorf("用户数组长度错误: expected %d, got %d", userCount, j.Get("users").Length())
	}

	if j.Get("users.0.profiles").Length() != profileCount {
		t.Errorf("配置数组长度错误: expected %d, got %d", profileCount, j.Get("users.0.profiles").Length())
	}

	t.Logf("✅ 大型对象数组测试通过：%d 个用户的完整嵌套结构", userCount)
}
