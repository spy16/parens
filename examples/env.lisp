; This program reads from environment variable ENV and checks if
; the value is "prod".

; checks if the value of ENV variable is prod.
(defn is-production ()
   (label current-env (env "ENV"))
   (== current-env "prod"))

; prints different message for different environments.
(cond
  ((is-production) (println "Production Environment!"))
  (true (println "Non-production environment")))
