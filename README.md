# vim-plugin-setup

`vim-plugin-setup` is a simple tool to help you setting up your .vimrc

## The reason I wrote this

1. One .vimrc is too complex for managing and setting for all vim plugins
2. Install plugins from a independent vimrc config file which can solve the deps of plugins and confs.
3. Run script will included in one vimrc conf would be better
4. Make stateful to avoid re-install plugins and re-run scripts

So you can found splited vimrc configurations inside vim-setups dir, open 'ycm.vimrc' for example:

```
" @require: github.com/Valloric/YouCompleteMe
"
" @run-script
" #!/bin/bash
" pushd ~/.vim/bundle/YouCompleteMe
" ./install.py  --clang-completer --gocode-completer --tern-completer
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
```

I add some `@requre` and `@run-script`/`@end-script` embedded in the comment, those scripts would be executed by `vim-plugin-setup` sequencially.

It's simple to combine vimrc conf and pthogen plugin manager, and now it works.

After run this, you will be automatically install some vim-plugins which I prefer to use:

- nerdtree
- ycm
- vim-go
- airline
- ...

## Install

```
go get github.com/karfield/vim-plugin-setup
```

## Usage

```
~/.go/bin/vim-plugin-setup install
```
