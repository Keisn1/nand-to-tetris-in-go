// function {{ .functionName }} {{ .nbrIterations }}
( {{ .functionName }} ){{range $index, $element := .Numbers}}

// Iteration #{{$index}}
// D=0 and *SP=D
@0
D=A
@SP
A=M
M=D

//SP++
@SP
M=M+1{{end}}
