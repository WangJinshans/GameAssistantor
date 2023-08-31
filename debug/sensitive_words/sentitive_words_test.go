package sensitive_words

import (
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	"regexp"
	"testing"
)


// HansCovertPinyin 中文汉字转拼音
func HansCovertPinyin(contents []string) []string {
	pinyinContents := make([]string, 0)
	for _, content := range contents {
		chineseReg := regexp.MustCompile("[\u4e00-\u9fa5]")
		if !chineseReg.Match([]byte(content)) {
			continue
		}

		// 只有中文才转
		pin := pinyin.New(content)
		pinStr, err := pin.Convert()
		fmt.Println(content, "->", pinStr)
		if err == nil {
			pinyinContents = append(pinyinContents, pinStr)
		}
	}
	return pinyinContents
}


// 前缀树匹配敏感词
func trieDemo(sensitiveWords []string, matchContents []string) {

	// 汉字转拼音
	pinyinContents := HansCovertPinyin(sensitiveWords)
	trie := NewSensitiveTrie()
	trie.AddWords(sensitiveWords)
	trie.AddWords(pinyinContents) // 添加拼音敏感词

	//trie.AddWords(pinyinContents)
	//for _, content := range contents {
	//	trie.AddWord(content)
	//}

	for _, srcText := range matchContents {
		matchSensitiveWords, replaceText := trie.Match(srcText)
		fmt.Println("srcText        -> ", srcText)
		fmt.Println("replaceText    -> ", replaceText)
		fmt.Println("sensitiveWords -> ", matchSensitiveWords)
		fmt.Println()
	}

	// 动态添加
	trie.AddWord("牛大大")
	content := "今天，牛大大挑战灰大大"
	matchSensitiveWords, replaceText := trie.Match(content)
	fmt.Println("srcText        -> ", content)
	fmt.Println("replaceText    -> ", replaceText)
	fmt.Println("sensitiveWords -> ", matchSensitiveWords)
}

func TestSensitiveTrie_Match(t *testing.T) {
	fmt.Println("--------- rune测试 ---------")

	sensitiveWords := []string{
		"傻逼",
		"傻叉",
		"垃圾",
		"妈的",
		"sb",
	}

	fmt.Println("\n--------- 汉字转拼音 ---------")
	pinyinContents := HansCovertPinyin(sensitiveWords)
	fmt.Println(pinyinContents)

	matchContents := []string{
		"你是一个大傻&逼，大傻 叉",
		"你是傻☺叉",
		"shabi东西",
		"他made东西",
		"什么垃圾打野，傻逼一样，叫你来开龙不来，SB",
		"正常的内容☺",
	}

	fmt.Println("\n--------- 前缀树匹配敏感词 ---------")
	trieDemo(sensitiveWords, matchContents)

}