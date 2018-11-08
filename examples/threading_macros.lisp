(defn square (n) (* n n))

; threading-first and threading-last macros can be
; used to reduce nesting.

; thread-first macro
(->
  (+ 1 2)
  (square)
  (println " is the square of the summation"))

; above expression without threading-first macro:
(println (square (+ 1 2)) " is the square of the summation")

; thread-last macro
(->>
  (+ 1 2)
  (square) ; since square takes only one argument, last/first are same
  (println "Square of the summation is "))

; above expression without threading-last macro:
(println "Square of the summation is " (square (+ 1 2)))