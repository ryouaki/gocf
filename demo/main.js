test("111111111", () => {
  
})

async function exec() {
  const ret = await Promise.all([new Promise((r) => {
    r(1)
  }), new Promise((r) => {
    r(2)
  })])
  console.log(ret)
  return ret
}

exec()