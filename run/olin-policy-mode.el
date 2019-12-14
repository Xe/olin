(require 'generic-x)

(define-generic-mode
  'olin-policy-mode                          ;; name of the mode
  '("##")                           ;; comments delimiter
  '("allow" "disallow" "ram-page-limit" "gas-limit")      ;; some keywords
  '()
  '("\\.policy$")                    ;; files that trigger this mode
  nil                              ;; any other functions to call
  "Olin Policy file mode"     ;; doc string
)
