# AnthropicAPI-go-wrapper

[![Go Reference](https://pkg.go.dev/badge/github.com/Stilltong/AnthropicAPI-go-wrapper/v2.svg)](https://pkg.go.dev/github.com/Stilltong/AnthropicAPI-go-wrapper/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/Stilltong/AnthropicAPI-go-wrapper/v2)](https://goreportcard.com/report/github.com/Stilltong/AnthropicAPI-go-wrapper/v2)
[![codecov](https://codecov.io/gh/Stilltong/AnthropicAPI-go-wrapper/graph/badge.svg?token=O6JSAOZORX)](https://codecov.io/gh/Stilltong/AnthropicAPI-go-wrapper)
[![Sanity check](https://github.com/Stilltong/AnthropicAPI-go-wrapper/actions/workflows/pr.yml/badge.svg)](https://github.com/Stilltong/AnthropicAPI-go-wrapper/actions/workflows/pr.yml)

This is an unofficial Anthropics Claude API wrapper for Go. It supports:

- Completions
- Streaming Completions
- Messages
- Streaming Messages
- Vision
- Tool use

## Installation

```
go get github.com/Stilltong/AnthropicAPI-go-wrapper/v2
```

The AnthropicAPI-go-wrapper requires Go version 1.21 or greater.

## Usage

### Messages example usage:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/Stilltong/AnthropicAPI-go-wrapper/v2"
)

func main() {	
//.....your codes
}
```

### Messages stream example usage:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/Stilltong/AnthropicAPI-go-wrapper/v2"
)

func main() {
//.....your codes
}
```

### Other examples:

<details>
<summary>Messages Vision example</summary>

```go
package main

import (
	"errors"
	"fmt"

	"github.com/Stilltong/AnthropicAPI-go-wrapper/v2"
)

func main() {
//.....your codes
}
```
</details>

<details>

<summary>Messages Tool use example</summary>

```go
package main

import (
	"context"
	"fmt"

	"github.com/Stilltong/AnthropicAPI-go-wrapper/v2"
	"github.com/Stilltong/AnthropicAPI-go-wrapper/v2/jsonschema"
)

func main() {
//.....your codes
}
```

</details>

## Acknowledgments
The following project had particular influence on AnthropicAPI-go-wrapper's design.

- [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai)