function exec() {
  ret = console&&console.log('aa', 11, {a:1})
  console.log(11, ret)

  ret = http.request('POST', 'https://wx.17u.cn/appapi/wxuser/login/2', {a:1}, {"code":"083s1r1w3U43WZ2vkv3w3uwC2T1s1r1O","scene":1001})
  console.log(JSON.stringify(ret), ret.data)
}
exec()