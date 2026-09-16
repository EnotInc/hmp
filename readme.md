# Help Me Please

## About

hmp is my second attempt on following the book "how to write an interpreter in go" by Thorsten Ball.

## Usage
You can run REPL by running `hmp`, or provide a file to run code from it:
```bash
hmp example\fibonacci.hmp
```

## Installation
```bash
git clone https://github.com/enotinc/hmp.git
cd hmt/cmd/hmp
go install # or with go build
```

## Sintax Example
```hmp
let age = 1;
age = 2;
let message = "help me please";
let res = 10 * (20/2);

const foo = ture;
foo = false; // => error


let arr = [1, 2, 3, 4];
let strc = {"name": "foo", "age": 21};

arr[1]; // => 2
strc["name"]; // => "foo"

let add = f(a, b) {
	return a + b;
};

add(1, 2) // => 3


// fibonacci example
let fibonacci = fn(n) {
	if (n < 2) { n }
	else { fibonacci(n-1) + fibonacci(n-2) }
}

fibonacci(12) // => 144

let count := 10;
for (count > 0) {
	print(count);
	count--;
}

```

> [!warning]
> buildin `print()` supports only alphabetic escape sequences (`\n`, `\t`, `\e` etc.).
> other sequences will be treated as separate bytes (`\033` = `{'\', '0', '3', '3'}`).
