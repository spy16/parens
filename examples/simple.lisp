;; need help ?
(doc println)
(dump-scope)

;; look under the hood
(inspect hello)
(println
  (type +))

;; calling variadic Go functions
(println "Hello" "from" "Parens!")

;; binding values
(label π 3.1412)

;; let starts a new scope
(let
  (label π 3)
  (printf "integer part of pi is %d\n" π))

;; value of π should now be reset to original
(printf "but real value of pi is %f\n" π)

;; let's use some cool looking characters since
;; parens supports UTF-8
(label ∂ 0)
(label ∑ +)
(label ≠
       (lambda [a b]
         (not
           (== a b))))

(println "1 + 2 ="
         (∑ 1 2))
(println "false and nil are not same ="
         (≠ false nil))

(label sign
       (lambda [in]
         (cond
           ((> in 0) "positive")
           ((< in 0) "negative")
           (true "zero"))))

;; defining a lambda
(label square
       (lambda [a]
         (* a a)))

;; calling a lambda, obviously
(printf "square of 2 is = %f\n"
        (square 2))

;; we need to do some math obviously
(printf "complex math answer %f\n"
        (* 1
           (- 2
              (+ 1
                 (/ 3 3)))))

;; time for real stuff.. fibonacci!
(defn fib [n]
      (cond
        ((< n 2) n)
        (true (+ (fib (- n 1)) (fib (- n 2))))))

;; what is the 10th number in the fibonacci sequence
(printf "10th number in the fibonacci sequence = %f\n"
        (fib 10))

;; parens supports character literals
(printf "Is \\u00A5 same as '¥'? = %t\n" (== \u00A5 \¥))

(println "Different Numbers: ")

(println
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
 )
