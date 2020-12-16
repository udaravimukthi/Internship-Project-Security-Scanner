# Process CPU infomation
sysdig -c topprocs_cpu
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/scan?command=sysdig%20-c%20topprocs_cpu

# Process network infomation
sysdig -c topprocs_net
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/scan?command=sysdig%20-c%20topprocs_net

# ip information
sysdig -i spy_ip
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/scan?command=sysdig%20-i%20spy_ip

# read files slower than 1 second
sysdig -c fileslower 1000
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/scan?command=sysdig%20-c%20fileslower%201000

# monitor network connections (system state)
sysdig -c netstat
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/scan?command=sysdig%20-c%20netstat

# HTTP requests log
sysdig -c httplog
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/scan?command=sysdig%20-c%20httplog

# Launched Programs
http://securityscanner.app.rancher.damian.k8s.ifsdevworld.com:31100/sysdigScan?command=sysdig%20-n%201000%20-p%22%25user.name%29%20%25evt.arg.path%22%20%22evt.type%3Dchdir%22






