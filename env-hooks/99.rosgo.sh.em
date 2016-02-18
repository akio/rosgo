@[if DEVELSPACE]@
if [ $GOPATH ]; then
    export GOPATH="@(CATKIN_DEVEL_PREFIX)/lib/go":"$GOPATH"
else
    export GOPATH="@(CATKIN_DEVEL_PREFIX)/lib/go"
fi
@[else]@
if [ $GOPATH ]; then
    export GOPATH="${CATKIN_ENV_HOOK_WORKSPACE}/lib/go":"$GOPATH"
else
    export GOPATH="${CATKIN_ENV_HOOK_WORKSPACE}/lib/go"
fi
@[end if]@
