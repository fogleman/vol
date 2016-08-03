## Volume Rendering in Go

Voxels are subsampled and multiplied against a projection matrix. Multiple color windows can be applied.

http://i.imgur.com/JrN75I7.gifv

### How-To

    go get -u github.com/fogleman/vol

Convert DICOMs to PNG using ImageMagick.

    $ mkdir png
    $ mogrify -path png/ -format png *.dcm

Run this script.

    $ vol png/
