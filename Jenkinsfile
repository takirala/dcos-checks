#!/usr/bin/env groovy

def master_branches = ["master", ] as String[]

if (master_branches.contains(env.BRANCH_NAME)) {
    // Rebuild main branch once a day
    properties([
        pipelineTriggers([cron('H H * * *')])
    ])
}

node('mesos-ubuntu') {
    stage ('checkout-scm') {
        dir("dcos-checks") {
            checkout scm
        }
        stash includes: 'dcos-checks/**', name: 'dcos-checks'
    }
    stage ('run-tests') {
		unstash 'dcos-checks'
        dir("dcos-checks") {
			sh 'make test'
        }
    }
}
