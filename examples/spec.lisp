; All the constructs supported by the Reader (with default read
; table.)

; strings --------------------------------------------------------
"Hello"                                                     ; simple string
"Hello\tWorld"                                              ; string with escape sequences

"this is going
to be a multi
line string"                                                ; a multi-line string

; characters -----------------------------------------------------
\a                                                          ; simple character literal representing 'a'
\newline                                                    ; special character literal representing '\n'
\u00A5                                                      ; unicode character literal representing '¥'

; numbers --------------------------------------------------------
1234                                                        ; simple integer (int64)
3.142                                                       ; simple double precision floating point (float64)
-1.23445                                                    ; negative floating point number
010                                                         ; octal representation of 8
-010                                                        ; octal representation of -8
0xF                                                         ; hexadecimal representation of 15
-0xAF                                                       ; negative hexadecimal number -175
1e3                                                         ; scientific notation for 1x10^3 = 1000
1.5e3                                                       ; scientific notation for 1.5x10^3 = 1500
10e-1                                                       ; scientific notation with negative exponent

; keywords -------------------------------------------------------
:key                                                        ; simple ASCII keyword
:find-Ψ                                                     ; a keyword with non non-ASCII

; symbols --------------------------------------------------------
hello                                                       ; simple ASCII symbol
calculate-λ                                                 ; symbol with non-ASCII

; quote/unquote --------------------------------------------------
'hello                                                      ; quoted symbol
'()                                                         ; quoted list
~(x 1)                                                      ; unquoting

; lists ----------------------------------------------------------
()                                                          ; empty list
(+ 1 2)                                                     ; function/macro/special form invocation
(+, 1, 2)                                                   ; same as above, forms separated by ","

; vectors --------------------------------------------------------
[]                                                          ; empty vector
[1 2 3 4]                                                   ; vector entries separated by space
[1, 2, 3, 4]                                                ; vector entries can be separated by "," as well
