; This example shows definition of a recursive function
(defn factorial [n]
      (cond
        ((== n 1) n)
        (true (* (factorial (- n 1)) n))))


; finally call the factorial function
(printf "10! = %f\n" (factorial 10))
