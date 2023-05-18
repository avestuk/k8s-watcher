# k8s-watcher

This is a implementation of watching a K8s ConfigMap mounted into a pod
intended to highlight the behaviour noted in the [blog post by Ahmet Alp
Balkan](https://ahmet.im/blog/kubernetes-inotify/)

The behaviour of the inotify library I have experience using does not
restablish watches once the configMap file has been removed. Instead you must
watch the parent directory and when events are triggered, re-read the file. 

If you run the deployments in a kubernetes cluster and edit the contents of the
configMap, you'll see that the direct-to-file watcher does not see the contents
change. 
