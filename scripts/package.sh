pushd out
for f in *; do
    [ -d "${f}" ] || continue
    cp ../LICENSE ../README.md $f
    zip -r $f.zip $f
done;
popd
