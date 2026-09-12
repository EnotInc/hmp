# Help Me Plese

hmp is a second attempt on following the book "how to write an interpreter in go" by Thorsten Ball

## Sintax

```hmp
let age = 1;
let message = "help me please";
let res = 10 * (20/2);

let arr = [1, 2, 3, 4];
let strc = {"name": "foo", "age": 21};

arr[1]; // => 2
strc["name"]; // => "foo"

let add = f(a, b) {
	return a + b;
};

add(1, 2)


// fibonacci example
let fibonacci = fn(x) {
	if (x == 0) {
		return 0;
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x-1) + fibonacci(x-2);
		}
	}
}
```
