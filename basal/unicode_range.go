package basal

// 汉字
var RuneNumber = RuneRange{'0', '9'}                                    //数字
var RuneUpper = RuneRange{'A', 'Z'}                                     //大写
var RuneLower = RuneRange{'a', 'z'}                                     //小写
var RuneChineseBasic = RuneRange{'\u4E00', '\u9FA5'}                    //基本汉字
var RuneChineseBasicEx = RuneRange{'\u9FA6', '\u9FFF'}                  //基本汉字补充
var RuneChineseExA = RuneRange{'\u3400', '\u4DBF'}                      //汉字扩展A
var RuneChineseExB = RuneRange{ToRune("\\u20000"), ToRune("\\u2A6DF")}  //汉字扩展B
var RuneChineseExC = RuneRange{ToRune("\\u2A700"), ToRune("\\u2B738")}  //汉字扩展C
var RuneChineseExD = RuneRange{ToRune("\\u2B740"), ToRune("\\u2B81D")}  //汉字扩展D
var RuneChineseExE = RuneRange{ToRune("\\u2B820"), ToRune("\\u2CEA1")}  //汉字扩展E
var RuneChineseExF = RuneRange{ToRune("\\u2CEB0"), ToRune("\\u2EBE0")}  //汉字扩展F
var RuneChineseExG = RuneRange{ToRune("\\u30000"), ToRune("\\u3134A")}  //汉字扩展G
var RuneChineseKx = RuneRange{'\u2F00', '\u2FD5'}                       //康熙部首
var RuneChineseKxEx = RuneRange{'\u2E80', '\u2EF3'}                     //部首扩展
var RuneChineseJr = RuneRange{'\uF900', '\uFAD9'}                       //兼容汉字
var RuneChineseJrEx = RuneRange{ToRune("\\u2F800"), ToRune("\\u2FA1D")} //兼容扩展
var RuneChinesePua = RuneRange{'\uE815', '\uE86F'}                      //PUA(GBK)部件
var RuneChinesePuaEx = RuneRange{'\uE400', '\uE5E8'}                    //部件扩展
var RuneChinesePuaSup = RuneRange{'\uE600', '\uE6CF'}                   //PUA增补
var RuneChineseBh = RuneRange{'\u31C0', '\u31E3'}                       //汉字笔画
var RuneChineseJg = RuneRange{'\u2FF0', '\u2FFB'}                       //汉字结构
var RuneChineseZy = RuneRange{'\u3105', '\u312F'}                       //汉语注音
var RuneChineseZyEx = RuneRange{'\u31A0', '\u31BA'}                     //注音扩展
var RuneChineseYq = RuneRange{'\u3007', '\u3007'}                       //〇

// 日语
var RuneJapaneseHir = RuneRange{'\u3040', '\u309F'} //平假名
var RuneJapaneseKat = RuneRange{'\u30A0', '\u30FF'} //片假名
var RuneJapaneseKPE = RuneRange{'\u31F0', '\u31FF'} //片假名音标扩展集

// 韩语
var RuneKoreanPSC = RuneRange{'\uAC00', '\uD7AF'} //谚文音节字符
var RuneKoreanHJ = RuneRange{'\u1100', '\u11FF'}  //谚文字母
var RuneKoreanPCL = RuneRange{'\u3130', '\u318F'} //谚文相容字母

// 中日韩
var RuneCJKA = RuneRange{'\uFF00', '\uFFEF'} //全角ASCII、全角中英文标点、半宽片假名、半宽平假名、半宽韩文字母
