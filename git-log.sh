#!/bin/bash

COUNT=${1:-1}

function markdown_1 {
  git log --pretty=format:'%s|%b' -$COUNT |
    while IFS='|' read title body
    do
      echo "* ${title}"
      if [[ "$body" != "" ]] ; then
        echo '      ```'
        echo "${body}" | sed 's/^/      /' | sed 's/CRLF/\n/g'
        echo '      ```'
      fi
  #    printf '* %s\n      ```\n%-20s\n      ```' "$title" "$body"
    done
}

function markdown_2() {
    git log -$COUNT --pretty=format:'* %s%n      ```%n%w(0,6,6)%b%n```%w(0,0,0)%n' | sed 's/\n\\s*```\\s*```\\s*\n/-/g'
}

function markdown_3() {
    git log -$COUNT --pretty=format:'* %s%n      ```%n%w(0,6,6)%b%n```%w(0,0,0)%n' | sed  '/[ \t\n]*```[ \t\n]*/p-/'
}

function markdown_4() {
    git log -$COUNT --pretty=tformat:'* %s%n      ```%n%w(0,6,6)%b%n```%w(0,0,0)%n' |
      sed '/[ \t\n]*```$/{
             $!{ N        # append the next line when not on the last line
               s/[ \t\n\r]*```[ \t\n\r\<00011\>]*/-/
                          # now test for a successful substitution, otherwise
                          #+  unpaired "a test" lines would be mis-handled
               t sub-yes  # branch_on_substitute (goto label :sub-yes)
               :sub-not   # a label (not essential; here to self document)
                          # if no substituion, print only the first line
               P          # pattern_first_line_print
               D          # pattern_ltrunc(line+nl)_top/cycle
               :sub-yes   # a label (the goto target of the 't' branch)
                          # fall through to final auto-pattern_print (2 lines)
             }
           }'
}

function markdown_5() {
    git log -$COUNT --pretty=tformat:'* %s%n      ```%n%w(0,6,6)%b```%w(0,0,0)%n' |
      sed '/```$/{
             $!{ N        # append the next line when not on the last line
               s/[ \t\n\r]*```[^a-zA-Z0-9]*```[^a-zA-Z0-9]*//
                          # now test for a successful substitution, otherwise
                          #+  unpaired "a test" lines would be mis-handled
               t sub-yes  # branch_on_substitute (goto label :sub-yes)
               :sub-not   # a label (not essential; here to self document)
                          # if no substituion, print only the first line
               P          # pattern_first_line_print
               D          # pattern_ltrunc(line+nl)_top/cycle
               :sub-yes   # a label (the goto target of the 't' branch)
                          # fall through to final auto-pattern_print (2 lines)
             }
           }'
}

function markdown_6() {
    git log -$COUNT --pretty=tformat:'* %s%n      %n%w(0,6,6)%b%w(0,0,0)%n' | sed -e 's/[ (]*\(+next\|+pr\)[ )]*/ /'
#    |
#      sed '/```$/{
#             $!{ N        # append the next line when not on the last line
#               s/[ \t\n\r]*```[^a-zA-Z0-9]*```[^a-zA-Z0-9]*//
#                          # now test for a successful substitution, otherwise
#                          #+  unpaired "a test" lines would be mis-handled
#               t sub-yes  # branch_on_substitute (goto label :sub-yes)
#               :sub-not   # a label (not essential; here to self document)
#                          # if no substituion, print only the first line
#               P          # pattern_first_line_print
#               D          # pattern_ltrunc(line+nl)_top/cycle
#               :sub-yes   # a label (the goto target of the 't' branch)
#                          # fall through to final auto-pattern_print (2 lines)
#             }
#           }'
}

markdown_6
