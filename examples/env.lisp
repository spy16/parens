(do
  ; This program reads from environment variable ENV and checks if
  ; the value is "prod".
  (defn is-production ()
     (label current-env (env "ENV"))
     (== current-env "prod"))
  (cond
    ((is-production) (println "Production Environment!"))
    (true (println "Non-production environment"))))
