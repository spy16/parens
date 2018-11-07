(do
	(defn square (n) (* n n))

	; threading macro and threading last macro can be
	; used to reduce nesting.
	(->>
		(-> (+ 1 2) (square))
		(println "Square of summation is ")))
