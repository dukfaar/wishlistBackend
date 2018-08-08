node {
    checkout scm
        
    stage('Docker Build') {
        docker.build('dukfaar/wishlistbackend')
    }

    stage('Update Service') {
        sh 'docker service update --force wishlistbackend_wishlistbackend'
    }
}
