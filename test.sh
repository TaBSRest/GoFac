#!/bin/bash
while getopts v:p:t:Show-No-Test-Files-Warning flag
do
    case "${flag}" in
        v) verbosity=${OPTARG};;
        p) package=${OPTARG};;
        t) target=${OPTARG};;
        # Show-No-Test-Files-Warning) warning=${OPTARG};;
    esac
done
echo "Verbosity: $verbosity";
echo "Package: $package";
echo "Target: $target";
# echo "Show-No-Test-Files-Warning" $warning";