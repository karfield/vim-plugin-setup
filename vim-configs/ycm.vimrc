" @require: github.com/Valloric/YouCompleteMe
"
" @run-script
" #!/bin/bash
" pushd ~/.vim/bundle/YouCompleteMe
" python2 ./install.py  --clang-completer --gocode-completer --tern-completer
" popd
" @end-script
"

" YCM settings
let g:ycm_key_list_select_completion = ['<tab>', '<c-p>']
let g:ycm_key_list_previous_completion = ['<shift-tab>', '<c-s-p>']
let g:ycm_key_invoke_completion = '<C-Space>'

" Trigger configuration. Do not use <tab> if you use
" https://github.com/Valloric/YouCompleteMe.
" let g:UltiSnipsExpandTrigger="<tab>"
let g:UltiSnipsJumpForwardTrigger="<c-b>"
let g:UltiSnipsJumpBackwardTrigger="<c-z>"

" If you want :UltiSnipsEdit to split your window.
let g:UltiSnipsEditSplit="vertical"
