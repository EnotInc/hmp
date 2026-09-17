# Help Me Please

## About

hmp is my second attempt on following the book "how to write an interpreter in go" by Thorsten Ball.

## Installation
```bash
git clone https://github.com/enotinc/hmp.git
cd hmp/cmd/hmp
go install # or with go build
```

## Usage
### Examples

You can fild a few examples of hmp scripts int `./example/` folder.

### Getting started

Let's assume that you already installed a `hmp` interpreter. To run your first script, you should provide a path to script file with `.hmp` ext.
```shell
hmp foo/bar.hmp
```
In this example we running `bar.hmp` file in `foo` dir.
Let's write our first script!. Instead of creating file we can use REPL. run `hmp` and you will start it.

Anyway, let's write hello world!
```js
print("hello world!") // 'hello world!'
```
And this is it! Our first script!
To quit REPL press `<ctrl+c>`

### Syntax
#### Variables
```js
const _my_constant = 5; // this value can't be changed
let _my_var = "hello";  //this one - can

_my_constant = 6; // => error
_my_var = 123;   // ok, now _my_var is an integer
```

##### Strings
```js
let s = "new"
s = s + " text"
let f = "another text"

let line = s + " " + f
print(line) // 'new text another text'

print(s[0]) // symbol - 'n'
```

##### Arrays
```js
let arr = [1, 2, 3]
print(arr[0]) // 1

arr = push(arr, "line")
print(arr) // [1, 2, 3, "line"], yes, arrays can contain different type of values
```

##### Hashmaps
```js
let m = {"name": "foo", "age": 21}
print(m) // {"name": "foo", "age": 21}

m["job"] = "jobless"
print(m) //{"name": "foo", "age": 21, "job": "jobless"}

m["job"] = "qa engineer"
print(m) //{"name": "foo", "age": 21, "job": "qa engineer"}
```

##### Boolean
```js
// boolean
let t = true
let f = false

let n = t || f // true
let m = t && f // false
let o = !t // false
let p = !f // true
```

#### Functions
```js
let add = fn(x,y) {
	return x+y;
}  // you can also type only x+y, without return keyword

let res = add(5, 10)
print(res) // 15
```

#### Control flow
```js
let dev = fn(x,y) {
	if (y==0) {
		return;
	} // if you wanna return from function, you should use return keyword with ';' semicolon symbol. This way function will ruturn null
	
	x/y; // just as I said earlier, this will return x/y, even without return keyword
}

const foo = dev(20, 5)
const bar = dev(20, 0)

print(foo, bar) // 4, null
```

```js
// guess.hmp
const val = 55;
let guess = fn(x) {
	if (x == val) {
		return "you guessed it!"
	} else {
		
		if (x < val) {
			return "value is greater"
		} else {
			return "value is less"
		}
	}
}

let res = guess(123)
print(res) // value is less
```
> [!warning]
> for now hmp doesn't support `else if (cond) {}` statement

#### Loops
let's continue on working our 'guess' game
```js
// guess.hmp
// [...]

let input = ""
let n = ""
for (n != val) {
	input = scanln() // getting user input
	n = atoi(input)  // parsing it to INTEGER
	
	const res = guess(n)
	print(res)
}
```

In hmp there are no 'while' loops, only 'for' loops
```js
for (cond) {}
```

if you want to run your loop forever, you can do this:
```js
for {
	// this will run forever
}
```

### Build-in functions
1. `atoi(<STRING>) <INTEGER>` - Ascii TO Integer, gets string as input and returns integer if parsed successfullly
2. `len(<ARRAY>|<STRING>) <INTEGER>` - gets array or a string and returns its length
3. `last(<ARRAY>) OBJ` - returns last element of given array
4. `tail(<ARRAY>) ARRAY` - returns new array, without first element
5. `push(<ARRAY>, <OBJ>) <ARRAY>` - returns new array, that combines given one and a new object
6. `pop(<ARRAY>) <ARRAY>` - returns new array without last element
7. `rand(<INTEGER>, <INTEGER>) <INTEGER>` - returns random value from min (first arg) and max (last arg)
8. `scanln() <STRING>` - returns user input from terminal
9. `print([OBJ])` - prints given list of objects
10. `exit([OBJ])` - prints given list of objects, and ends the program
11. `exists(<STRING>) <BOOLEAN>` - checks if entry (file or dir) exists on given path
12. `create(<STRING>)` - creates file with provided name
13. `read(<STRING>) <STRING>` - reads given file data and returns it
14. `write(<STRING>, <STRING>)` - writes given data (2nd arg) on file (1st arg)
15. `delete(<STRING>)` - deletes entry on provided path
16. `rename(<STRING>, <STRING>)` - renames entry from old to new one
17. `args() [<STRING>]` - returns an ARRAY of strings, with provided args on running script (1-st arg is script filename)

> [!warning]
> buildin `print()` supports only alphabetic escape sequences (`\n`, `\t`, `\e` etc.).
> other sequences will be treated as separate bytes (`\033` = `{'\', '0', '3', '3'}`).
