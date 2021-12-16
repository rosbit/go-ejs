console.log("env TZ =>", js.env("TZ"))
console.log("env age =>", js.env("age"))
console.log("call golang add(10, 2) =>", add(10, 2))
console.log("formatTime(1639632557) =>", formatTimestamp(1639632557))
// js.call("dump", 1, 2)
// js.call("dump1", 1, "hello")

function test(a1, b1) {
	console.log("---- test ----")
	var pi = Math.round(3.1415926*100)/100
	console.log(pi)
	console.log(a.name, a.age)
	console.log(m.a[0])
	return {a: a1, b: b1}
}
