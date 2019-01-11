#!/bin/bash

function prompt() {
    echo -n -e "\033[1;32m?\033[0m \033[1m$1\033[0m ($2) "
}

function replace() {
    sed -E -e "$1" $2 > $2.new
    mv -f $2.new $2
}

originalPackageName="github.com/sagikazarmark/modern-go-application"
originalBinaryName="modern-go-application"

defaultPackageName=${PWD##*src/}
prompt "Package name" ${defaultPackageName}
read packageName
packageName=$(echo "${packageName:-${defaultPackageName}}" | sed 's/[[:space:]]//g')

defaultProjectName=$(basename ${packageName})
prompt "Project name" ${defaultProjectName}
read projectName
projectName=$(echo "${projectName:-${defaultProjectName}}" | sed 's/[[:space:]]//g')

prompt "Binary name" ${projectName}
read binaryName
binaryName=$(echo "${binaryName:-${projectName}}" | sed 's/[[:space:]]//g')

prompt "Service name" ${projectName}
read serviceName
serviceName=$(echo "${serviceName:-${projectName}}" | sed 's/[[:space:]]//g')

defaultFriendlyServiceName=$(echo "${serviceName}" | sed -e 's/-/ /g;' | awk '{for(i=1;i<=NF;i++){ $i=toupper(substr($i,1,1)) substr($i,2) }}1')
prompt "Friendly service name" "${defaultFriendlyServiceName}"
read friendlyServiceName
friendlyServiceName=${friendlyServiceName:-${defaultFriendlyServiceName}}

prompt "Remove init script" "y/N"
read removeInit
removeInit=${removeInit:-n}

# IDE configuration
mv .idea/project.iml .idea/${projectName}.iml
replace 's|.idea/project.iml|.idea/'${projectName}'.iml|g' .idea/modules.xml

# Run configurations
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/All_tests.xml
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/Debug.xml
replace 's|value="\$PROJECT_DIR\$\/cmd\/modern-go-application\/"|value="$PROJECT_DIR$/cmd/'${binaryName}'/"|' .idea/runConfigurations/Debug.xml
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/Integration_tests.xml
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/Tests.xml

# Binary name
mv cmd/${originalBinaryName} cmd/${binaryName}

# Variables
replace 's|ServiceName = "mga"|ServiceName = "'${serviceName}'"|' cmd/${binaryName}/vars.go
replace 's|FriendlyServiceName = "Modern Go Application"|FriendlyServiceName = "'"${friendlyServiceName}"'"|' cmd/${binaryName}/vars.go

# Makefile
replace "s|^PACKAGE = .*|PACKAGE = ${packageName}|" Makefile
replace "s|^BUILD_PACKAGE \??= .*|BUILD_PACKAGE = \${PACKAGE}/cmd/${binaryName}|" Makefile
replace "s|^BINARY_NAME \?= .*|BINARY_NAME \?= ${binaryName}|" Makefile

# Other project files
declare -a files=(".circleci/config.yml" ".gitlab-ci.yml" "CHANGELOG.md" "Dockerfile")
for file in "${files[@]}"; do
    if [[ -f "${file}" ]]; then
        replace "s|${originalPackageName}|${packageName}|" ${file}
    fi
done

# Example code
find cmd/ -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done
find internal/ -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done

if [[ "${removeInit}" != "n" && "${removeInit}" != "N" ]]; then
    rm "$0"
fi
