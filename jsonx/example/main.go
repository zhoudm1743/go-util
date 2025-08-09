package main

import (
	"fmt"

	"github.com/zhoudm1743/go-util/jsonx"
)

func main() {
	fmt.Println("🚀 JSONx 包使用示例")
	fmt.Println("===================")

	// 1. 基础操作示例
	fmt.Println("\n1. 基础操作示例")
	basicExample()

	// 2. 链式调用示例
	fmt.Println("\n2. 链式调用示例")
	chainExample()

	// 3. 构建器模式示例
	fmt.Println("\n3. 构建器模式示例")
	builderExample()

	// 4. 路径操作示例
	fmt.Println("\n4. 路径操作示例")
	pathExample()

	// 5. 数组操作示例
	fmt.Println("\n5. 数组操作示例")
	arrayExample()

	// 6. 高级功能示例
	fmt.Println("\n6. 高级功能示例")
	advancedExample()

	// 7. 实用工具示例
	fmt.Println("\n7. 实用工具示例")
	utilityExample()
}

func basicExample() {
	// 解析 JSON 字符串
	jsonStr := `{"name": "张三", "age": 25, "city": "北京"}`
	j := jsonx.Parse(jsonStr)

	// 获取值
	name := j.Get("name").String()
	age := j.Get("age").Int()
	fmt.Printf("姓名: %s, 年龄: %d\n", name, age)

	// 设置值
	j.Set("email", "zhangsan@example.com")
	j.Set("age", 26)

	// 输出修改后的 JSON
	result, _ := j.ToPrettyJSON()
	fmt.Printf("修改后的 JSON:\n%s\n", result)
}

func chainExample() {
	// 创建 JSON 对象并链式操作
	result := jsonx.Object().
		Set("user.name", "李四").
		Set("user.profile.age", 30).
		Set("user.profile.skills", []interface{}{"Go", "Python", "JavaScript"}).
		Set("user.active", true).
		Set("timestamp", 1640995200)

	// 链式获取和转换
	userName := result.Get("user.name").String()
	userAge := result.Get("user.profile.age").Int()
	isActive := result.Get("user.active").Bool()

	fmt.Printf("用户: %s, 年龄: %d, 活跃: %t\n", userName, userAge, isActive)

	// 输出完整 JSON
	jsonStr, _ := result.ToPrettyJSON()
	fmt.Printf("完整 JSON:\n%s\n", jsonStr)
}

func builderExample() {
	// 使用构建器创建复杂对象
	user := jsonx.NewBuilder().
		AddString("name", "王五").
		AddInt("age", 28).
		AddBool("verified", true).
		AddObject("address", jsonx.NewBuilder().
			AddString("street", "长安街").
			AddString("city", "北京").
			AddString("zipcode", "100000").
			Build()).
		AddArray("hobbies", jsonx.QuickArray("读书", "游泳", "编程")).
		Build()

	fmt.Printf("构建器创建的用户对象:\n")
	result, _ := user.ToPrettyJSON()
	fmt.Println(result)

	// 使用快速构建器
	product := jsonx.QuickObject(map[string]interface{}{
		"id":    "P001",
		"name":  "智能手机",
		"price": 2999.99,
		"tags":  []interface{}{"电子产品", "手机", "智能设备"},
	})

	fmt.Printf("快速构建的产品对象:\n")
	productJSON, _ := product.ToPrettyJSON()
	fmt.Println(productJSON)
}

func pathExample() {
	// 复杂的嵌套 JSON 结构
	jsonStr := `{
		"company": {
			"name": "科技有限公司",
			"departments": [
				{
					"name": "研发部",
					"employees": [
						{"name": "张三", "position": "工程师"},
						{"name": "李四", "position": "架构师"}
					]
				},
				{
					"name": "市场部",
					"employees": [
						{"name": "王五", "position": "经理"}
					]
				}
			]
		}
	}`

	j := jsonx.Parse(jsonStr)

	// 深度路径访问
	companyName := j.Get("company.name").String()
	firstDeptName := j.Get("company.departments.0.name").String()
	firstEmployee := j.Get("company.departments.0.employees.0.name").String()

	fmt.Printf("公司: %s\n", companyName)
	fmt.Printf("第一个部门: %s\n", firstDeptName)
	fmt.Printf("第一个员工: %s\n", firstEmployee)

	// 修改深度嵌套的值
	j.Set("company.departments.0.employees.0.salary", 8000)
	j.Set("company.founded", "2020")

	// 检查路径是否存在
	fmt.Printf("公司成立时间存在: %t\n", j.Has("company.founded"))
	fmt.Printf("CEO 信息存在: %t\n", j.Has("company.ceo"))

	// 删除路径
	j.Delete("company.departments.1")
	fmt.Printf("删除市场部后的结构:\n")
	result, _ := j.ToPrettyJSON()
	fmt.Println(result)
}

func arrayExample() {
	// 创建和操作数组
	arr := jsonx.Array().
		Append("第一项").
		Append(42).
		Append(true).
		Append(map[string]interface{}{"key": "value"})

	fmt.Printf("数组长度: %d\n", arr.Length())
	fmt.Printf("第二个元素: %d\n", arr.Index(1).Int())

	// 数组迭代
	fmt.Println("遍历数组:")
	arr.ForEach(func(key string, value *jsonx.JSON) bool {
		fmt.Printf("  [%s] = %s (%s)\n", key, value.String(), jsonx.GetType(value))
		return true // 继续遍历
	})

	// 数组映射
	numbers := jsonx.QuickArray(1, 2, 3, 4, 5)
	doubled := numbers.Map(func(key string, value *jsonx.JSON) interface{} {
		return value.Int() * 2
	})

	fmt.Printf("原数组: ")
	numsJSON, _ := numbers.ToJSON()
	fmt.Println(numsJSON)

	fmt.Printf("翻倍后: ")
	doubledJSON, _ := doubled.ToJSON()
	fmt.Println(doubledJSON)

	// 数组过滤
	filtered := numbers.Filter(func(key string, value *jsonx.JSON) bool {
		return value.Int()%2 == 0 // 只保留偶数
	})

	fmt.Printf("过滤偶数: ")
	filteredJSON, _ := filtered.ToJSON()
	fmt.Println(filteredJSON)
}

func advancedExample() {
	// 克隆和合并
	obj1 := jsonx.QuickObject(map[string]interface{}{
		"name":  "产品A",
		"price": 100,
		"features": map[string]interface{}{
			"color": "red",
			"size":  "large",
		},
	})

	obj2 := jsonx.QuickObject(map[string]interface{}{
		"price":    120, // 覆盖价格
		"category": "电子产品",
		"features": map[string]interface{}{
			"weight": "1kg",
			"color":  "blue", // 覆盖颜色
		},
	})

	// 浅合并
	merged := obj1.Clone().Merge(obj2)
	fmt.Println("浅合并结果:")
	mergedJSON, _ := merged.ToPrettyJSON()
	fmt.Println(mergedJSON)

	// 深度合并
	deepMerged := obj1.Clone().DeepMerge(obj2)
	fmt.Println("深度合并结果:")
	deepMergedJSON, _ := deepMerged.ToPrettyJSON()
	fmt.Println(deepMergedJSON)

	// 结构转换
	user := struct {
		Name   string   `json:"name"`
		Age    int      `json:"age"`
		Tags   []string `json:"tags"`
		Active bool     `json:"active"`
	}{
		Name:   "测试用户",
		Age:    25,
		Tags:   []string{"开发者", "Go"},
		Active: true,
	}

	userJSON := jsonx.FromStruct(user)
	fmt.Println("从结构体创建的 JSON:")
	userJSONStr, _ := userJSON.ToPrettyJSON()
	fmt.Println(userJSONStr)
}

func utilityExample() {
	// 模板构建器
	template := `{
		"user": "{{username}}",
		"message": "{{message}}",
		"timestamp": {{timestamp}},
		"active": {{active}}
	}`

	templateJSON := jsonx.NewTemplate(template).
		Set("username", "模板用户").
		Set("message", "这是一条模板消息").
		Set("timestamp", 1640995200).
		Set("active", true).
		Build()

	fmt.Println("模板构建的 JSON:")
	templateResult, _ := templateJSON.ToPrettyJSON()
	fmt.Println(templateResult)

	// JSON 扁平化和反扁平化
	complex := jsonx.QuickObject(map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"name": "扁平化测试",
				"age":  30,
			},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "zh-CN",
			},
		},
		"items": []interface{}{"item1", "item2"},
	})

	// 扁平化
	flattened := jsonx.Flatten(complex)
	fmt.Println("扁平化结果:")
	for k, v := range flattened {
		fmt.Printf("  %s = %v\n", k, v)
	}

	// 反扁平化
	unflattened := jsonx.Unflatten(flattened)
	fmt.Println("反扁平化结果:")
	unflattenedJSON, _ := unflattened.ToPrettyJSON()
	fmt.Println(unflattenedJSON)

	// 字段选择和排除
	original := jsonx.QuickObject(map[string]interface{}{
		"id":       1,
		"name":     "测试产品",
		"price":    99.99,
		"internal": "内部信息",
		"secret":   "机密数据",
	})

	// 只选择公开字段
	public := jsonx.Pick(original, "id", "name", "price")
	fmt.Println("公开字段:")
	publicJSON, _ := public.ToPrettyJSON()
	fmt.Println(publicJSON)

	// 排除敏感字段
	safe := jsonx.Omit(original, "internal", "secret")
	fmt.Println("安全字段:")
	safeJSON, _ := safe.ToPrettyJSON()
	fmt.Println(safeJSON)

	// 实用工具函数
	testJSON := jsonx.QuickObject(map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "深度嵌套",
			},
		},
		"array": []interface{}{1, 2, 3},
	})

	fmt.Printf("JSON 大小: %d 字节\n", jsonx.Size(testJSON))
	fmt.Printf("JSON 深度: %d\n", jsonx.Depth(testJSON))
	fmt.Printf("JSON 类型: %s\n", jsonx.GetType(testJSON))

	// JSON 字符串验证和格式化
	messyJSON := `{"name":"test","age":25,"active":true}`
	fmt.Printf("原始 JSON: %s\n", messyJSON)
	fmt.Printf("是否有效: %t\n", jsonx.IsValid(messyJSON))

	prettyJSON, _ := jsonx.Pretty(messyJSON)
	fmt.Printf("格式化后:\n%s\n", prettyJSON)

	minifiedJSON, _ := jsonx.Minify(prettyJSON)
	fmt.Printf("压缩后: %s\n", minifiedJSON)
}
