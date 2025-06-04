package quickstart

import "testing"

func TestChat(t *testing.T) {
	role := "程序员心理健康顾问"
	style := "积极乐观"
	chatHistory := []string{
		"我最近在编程时感到很沮丧，怎么办？",
		"保持积极的心态是很重要的。你可以尝试以下方法：\n1. 定期休息，避免过度疲劳。\n2. 与同事交流，分享你的感受。\n3. 参加一些有趣的编程活动，放松心情。",
		"我觉得自己在编程上遇到了瓶颈，怎么办？",
		"遇到瓶颈是很常见的现象。你可以尝试以下方法：\n1. 换个角度思考问题，可能会有新的灵感。\n2. 学习新的技术或工具，拓宽你的视野。\n3. 与其他程序员交流，获取不同的观点和建议。",
	}
	question := "如何保持编程时的积极心态？"

	Chat(role, style, question, chatHistory)
}
