package main

import ("fmt"
		"net/http"
		"net/url"
		"os/exec"
		"time"
		"strings"
	)

func main() {
	static := http.FileServer(http.Dir("./static"))
 	http.Handle("/", static)

	scanresults := http.FileServer(http.Dir("/scanresults"))
	http.Handle("/scanresults/", http.StripPrefix("/scanresults/", scanresults))

	http.HandleFunc("/sysdigScan", sysdigScan_handler)

	http.HandleFunc("/trivyScan", trivyScan_handler)

	http.HandleFunc("/clairDeploy", clairDeploy_handler)
	http.HandleFunc("/clairDeployWithPersistence", clairDeployWithPersistence_handler)
	http.HandleFunc("/clairScan", clairScan_handler)

	http.HandleFunc("/csysdigScan", csysdigScan_handler)

	http.HandleFunc("/kubebenchScan", kubebenchScan_handler)

	http.HandleFunc("/kubehunterScan", kubehunterScan_handler)

	http.HandleFunc("/falcoDeploy", falcoDeploy_handler)

	http.HandleFunc("/fetchimages", fetchimages_handler)

	http.HandleFunc("/enablePersistence", enablePersistence_handler)

	http.ListenAndServe(":8080", nil)
}

func sysdigScan_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		deleteSysdigJobOut, deleteSysdigJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete job security-scanner-sysdig -n security-scanner").CombinedOutput() 
		if deleteSysdigJobOutErr != nil {
			fmt.Println(deleteSysdigJobOutErr)
		}
		fmt.Println(string(deleteSysdigJobOut))

		deleteSysdigResultsJobOut, deleteSysdigResultsJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && rm -rf /scanresults/capture.scap").CombinedOutput() 
		if deleteSysdigResultsJobOutErr != nil {
			fmt.Println(deleteSysdigResultsJobOutErr)
		}
		fmt.Println(string(deleteSysdigResultsJobOut))

		copySysdigJobOut, SysdigJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-sysdig-job.yaml /tmp/security-scanner-sysdig-job.yaml").CombinedOutput() 
		if SysdigJobOutErr != nil {
			fmt.Println(SysdigJobOutErr)
		}
		fmt.Println(string(copySysdigJobOut))

		fmt.Println(r.URL.Query().Get("command"))
		fmt.Println("sed -i 's/COMMAND/" + r.URL.Query().Get("command") + "/g' /tmp/security-scanner-sysdig-job.yaml")

		SetSysdigJobOut, SysdigJobSetErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's/COMMAND_PLACEHOLDER/" + r.URL.Query().Get("command") + "/g' /tmp/security-scanner-sysdig-job.yaml").CombinedOutput() 
		if SysdigJobSetErr != nil {
			fmt.Println(SysdigJobSetErr)
		}
		fmt.Println(string(SetSysdigJobOut))

		setScanSecondsSysdigJobOut, setScanSecondsSysdigJobSetErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's/SCANSECONDS_PLACEHOLDER/" + r.URL.Query().Get("scanseconds") + "/g' /tmp/security-scanner-sysdig-job.yaml").CombinedOutput() 
		if setScanSecondsSysdigJobSetErr != nil {
			fmt.Println(setScanSecondsSysdigJobSetErr)
		}
		fmt.Println(string(setScanSecondsSysdigJobOut))

		out, err := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-sysdig-job.yaml").CombinedOutput() 
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		
		podName := ""
		timeout := 5
		for {
			if timeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				podNameOut, podNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Running --selector=app=security-scanner-sysdig -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if podNameErr != nil {
					fmt.Println(podNameErr)
				}
				if len(podName) > 0 && strings.TrimSuffix(string(podNameOut[17:len(podNameOut)]), "\n") == podName {
					break
				} else {
					time.Sleep(20 * time.Second) 
					timeout--  
				}

				podName = string(podNameOut[17:len(podNameOut)])
				podName = strings.TrimSuffix(podName, "\n")	                     
			}
		}

		fmt.Println("Pod " + podName)
		
		succeededPodName := ""
		succeededTimeout := 120        
		for {
			if succeededTimeout <= 0 {
				break
			} else {
				if len(succeededPodName) > 0 {
					break
				} else {
					time.Sleep(1 * time.Second) 
					succeededTimeout--  
				} 

				permissionSysdigResultsJobOut, permissionSysdigResultsJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && mkdir -p /scanresults && chmod -R 777 /scanresults").CombinedOutput() 
				if permissionSysdigResultsJobOutErr != nil {
					fmt.Println(permissionSysdigResultsJobOutErr)
				}
				fmt.Println(string(permissionSysdigResultsJobOut))
		
				copySysdigResultOut, copySysdigResultErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl cp " + podName + ":/scanresults/security-scanner-sysdig-job-scan-results.txt /scanresults/security-scanner-sysdig-job-scan-results.txt -n security-scanner").CombinedOutput()
				if copySysdigResultErr != nil {
					fmt.Println(copySysdigResultErr)
				}
				fmt.Println(string(copySysdigResultOut))
		        
				copyCaptureFileSysdigResultOut, copyCaptureFileSysdigResultErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl cp " + podName + ":/scanresults/capture.scap /scanresults/capture.scap -n security-scanner").CombinedOutput()
				if copyCaptureFileSysdigResultErr != nil {
					fmt.Println(copyCaptureFileSysdigResultErr)
				}
				fmt.Println(string(copyCaptureFileSysdigResultOut))
		
				succeededPodNameOut, succeededPodNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Succeeded --selector=app=security-scanner-sysdig -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if succeededPodNameErr != nil {
					fmt.Println(succeededPodNameErr)
				}
				succeededPodName = string(succeededPodNameOut[17:len(succeededPodNameOut)])
				succeededPodName = strings.TrimSuffix(succeededPodName, "\n")	                     
			}
		}

		fmt.Println("Succeeded Pod " + succeededPodName)

		fmt.Fprintf(w, `<h1>Security Scanner - Sysdig</h1>
		<p><button><a href="scanresults/security-scanner-sysdig-job-scan-results.txt" target="_blank" style="text-decoration:none">Show Scan Results</a></button></p>
        <p><button><a href="scanresults/capture.scap" download="capture.scap">Download Scap File</a></button></p>
		<p><button onclick="window.open(window.location.origin.replace('31100', '31101/#/capture/%s/views/overview'));">Inspect with Sysdig Inspect</button></p>
		`, url.QueryEscape("/scanresults/capture.scap"))
	}
}

func trivyScan_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		deleteTrivyJobOut, deleteTrivyJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete job security-scanner-trivy -n security-scanner").CombinedOutput() 
		if deleteTrivyJobOutErr != nil {
			fmt.Println(deleteTrivyJobOutErr)
		}
		fmt.Println(string(deleteTrivyJobOut))

		deleteTrivyResultsJobOut, deleteTrivyResultsJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && rm -rf /scanresults/security-scanner-trivy-job-scan-results.txt").CombinedOutput() 
		if deleteTrivyResultsJobOutErr != nil {
			fmt.Println(deleteTrivyResultsJobOutErr)
		}
		fmt.Println(string(deleteTrivyResultsJobOut))

		copyTrivyJobOut, trivyJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-trivy-job.yaml /tmp/security-scanner-trivy-job.yaml").CombinedOutput() 
		if trivyJobOutErr != nil {
			fmt.Println(trivyJobOutErr)
		}
		fmt.Println(string(copyTrivyJobOut))

		fmt.Println(r.URL.Query().Get("command"))
		fmt.Println("sed -i 's~COMMAND~" + r.URL.Query().Get("command") + "~g' /tmp/security-scanner-trivy-job.yaml")

		trivyCommandOut, trivyCommandErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~COMMAND_PLACEHOLDER~" + r.URL.Query().Get("command") + "~g' /tmp/security-scanner-trivy-job.yaml").CombinedOutput() 
		if trivyCommandErr != nil {
			fmt.Println(trivyCommandErr)
		}
		fmt.Println(string(trivyCommandOut))

		registryUsername := string(r.URL.Query().Get("registryUsername"))

		registryUsernameOut, registryUsernameErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~REGISTRYUSERNAME_PLACEHOLDER~" + registryUsername + "~g' /tmp/security-scanner-trivy-job.yaml").CombinedOutput() 
		if registryUsernameErr != nil {
			fmt.Println(registryUsernameErr)
		}
		fmt.Println(string(registryUsernameOut))

		registryPassword := string(r.URL.Query().Get("registryPassword"))
	
		registryPasswordOut, registryPasswordErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~REGISTRYPASSWORD_PLACEHOLDER~" + registryPassword + "~g' /tmp/security-scanner-trivy-job.yaml").CombinedOutput() 
		if registryPasswordErr != nil {
			fmt.Println(registryPasswordErr)
		}
		fmt.Println(string(registryPasswordOut))

		outTrivy, errTrivy := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-trivy-job.yaml").CombinedOutput() 
		if errTrivy != nil {
			fmt.Println(errTrivy)
		}
		fmt.Println(string(outTrivy))

		podName := ""
		timeout := 10
		for {
			if timeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				if len(podName) > 0{
					break
				} else {
					time.Sleep(10 * time.Second) 
					timeout--  
				} 

				podNameOut, podNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Running --selector=app=security-scanner-trivy -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if podNameErr != nil {
					fmt.Println(podNameErr)
				}
				podName = string(podNameOut[17:len(podNameOut)])
				podName = strings.TrimSuffix(podName, "\n")	                     
			}
		}
		fmt.Println("Pod " + podName)

		succeededPodName := ""
		succeededTimeout := 300
		for {
			if succeededTimeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				if len(succeededPodName) > 0 {
					break
				} else {
					time.Sleep(1 * time.Second) 
					succeededTimeout--  
				} 

				copyTrivyResultOut, trivyJobResultErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl cp " + podName + ":/scanresults/security-scanner-trivy-job-scan-results.txt /scanresults/security-scanner-trivy-job-scan-results.txt -n security-scanner").CombinedOutput()
				if trivyJobResultErr != nil {
					fmt.Println(trivyJobResultErr)
				}
				fmt.Println(string(copyTrivyResultOut))
		
				succeededPodNameOut, succeededPodNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Succeeded --selector=app=security-scanner-trivy -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if succeededPodNameErr != nil {
					fmt.Println(succeededPodNameErr)
				}
				succeededPodName = string(succeededPodNameOut[17:len(succeededPodNameOut)])
				succeededPodName = strings.TrimSuffix(succeededPodName, "\n")	                     
			}
		}

		fmt.Println("Succeeded Pod " + succeededPodName)

		fmt.Fprintf(w,`<h1>Security Scanner - Trivy</h1>
		<p><button><a href="scanresults/security-scanner-trivy-job-scan-results.txt" target="_blank" style="text-decoration:none">Show Scan Results</a></button></p>
		`)	
	}
}


func clairDeploy_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		deletePostgresServiceOut, deletePostgresServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete service postgres-svc -n security-scanner").CombinedOutput() 
		if deletePostgresServiceErr != nil {
			fmt.Println(deletePostgresServiceErr)
		}
		fmt.Println(string(deletePostgresServiceOut))

		copyPostgresServiceOut, copyPostgresServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-postgres-service.yaml /tmp/security-scanner-clair-postgres-service.yaml").CombinedOutput() 
		if copyPostgresServiceErr != nil {
			fmt.Println(copyPostgresServiceErr)
		}
		fmt.Println(string(copyPostgresServiceOut))

		postgresServiceOut, postgresServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-postgres-service.yaml").CombinedOutput() 
		if postgresServiceErr != nil {
			fmt.Println(postgresServiceErr)
		}
		fmt.Println(string(postgresServiceOut))

		deletePostgresOut, deletePostgresErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete deployment security-scanner-clair-postgres -n security-scanner").CombinedOutput() 
		if deletePostgresErr != nil {
			fmt.Println(deletePostgresErr)
		}
		fmt.Println(string(deletePostgresOut))

		copyPostgresOut, copyPostgresErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-postgres-deployment.yaml /tmp/security-scanner-clair-postgres-deployment.yaml").CombinedOutput() 
		if copyPostgresErr != nil {
			fmt.Println(copyPostgresErr)
		}
		fmt.Println(string(copyPostgresOut))

		postgresOut, postgresErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-postgres-deployment.yaml").CombinedOutput() 
		if postgresErr != nil {
			fmt.Println(postgresErr)
		}
		fmt.Println(string(postgresOut))

		deleteClairServerOut, deleteClairServerErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete deployment security-scanner-clair-server -n security-scanner").CombinedOutput() 
		if deleteClairServerErr != nil {
			fmt.Println(deleteClairServerErr)
		}
		fmt.Println(string(deleteClairServerOut))

		time.Sleep(20 * time.Second)

		copyClairServerOut, copyClairServerErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-server-deployment.yaml /tmp/security-scanner-clair-server-deployment.yaml").CombinedOutput() 
		if copyClairServerErr != nil {
			fmt.Println(copyClairServerErr)
		}
		fmt.Println(string(copyClairServerOut))

		clairServerOut, clairServerOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-server-deployment.yaml").CombinedOutput() 
		if clairServerOutErr != nil {
			fmt.Println(clairServerOutErr)
		}
		fmt.Println(string(clairServerOut))
	}	
}

func clairDeployWithPersistence_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		
		copyPostgresPVCOut, copyPostgresPVCErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-postgres-pvc.yaml /tmp/security-scanner-clair-postgres-pvc.yaml").CombinedOutput() 
		if copyPostgresPVCErr != nil {
			fmt.Println(copyPostgresPVCErr)
		}
		fmt.Println(string(copyPostgresPVCOut))

		postgresPVCOut, postgresPVCErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-postgres-pvc.yaml").CombinedOutput() 
		if postgresPVCErr != nil {
			fmt.Println(postgresPVCErr)
		}
		fmt.Println(string(postgresPVCOut))

		deletePostgresServiceOut, deletePostgresServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete service postgres-svc -n security-scanner").CombinedOutput() 
		if deletePostgresServiceErr != nil {
			fmt.Println(deletePostgresServiceErr)
		}
		fmt.Println(string(deletePostgresServiceOut))

		copyPostgresServiceOut, copyPostgresServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-postgres-service.yaml /tmp/security-scanner-clair-postgres-service.yaml").CombinedOutput() 
		if copyPostgresServiceErr != nil {
			fmt.Println(copyPostgresServiceErr)
		}
		fmt.Println(string(copyPostgresServiceOut))

		postgresServiceOut, postgresServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-postgres-service.yaml").CombinedOutput() 
		if postgresServiceErr != nil {
			fmt.Println(postgresServiceErr)
		}
		fmt.Println(string(postgresServiceOut))

		deletePostgresOut, deletePostgresErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete deployment security-scanner-clair-postgres -n security-scanner").CombinedOutput() 
		if deletePostgresErr != nil {
			fmt.Println(deletePostgresErr)
		}
		fmt.Println(string(deletePostgresOut))

		copyPostgresOut, copyPostgresErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-postgres-deployment-pv.yaml /tmp/security-scanner-clair-postgres-deployment-pv.yaml").CombinedOutput() 
		if copyPostgresErr != nil {
			fmt.Println(copyPostgresErr)
		}
		fmt.Println(string(copyPostgresOut))

		postgresOut, postgresErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-postgres-deployment-pv.yaml").CombinedOutput() 
		if postgresErr != nil {
			fmt.Println(postgresErr)
		}
		fmt.Println(string(postgresOut))

		deleteClairServerOut, deleteClairServerErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete deployment security-scanner-clair-server -n security-scanner").CombinedOutput() 
		if deleteClairServerErr != nil {
			fmt.Println(deleteClairServerErr)
		}
		fmt.Println(string(deleteClairServerOut))

		time.Sleep(20 * time.Second)

		copyClairServerOut, copyClairServerErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-server-deployment.yaml /tmp/security-scanner-clair-server-deployment.yaml").CombinedOutput() 
		if copyClairServerErr != nil {
			fmt.Println(copyClairServerErr)
		}
		fmt.Println(string(copyClairServerOut))

		clairServerOut, clairServerOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-server-deployment.yaml").CombinedOutput() 
		if clairServerOutErr != nil {
			fmt.Println(clairServerOutErr)
		}
		fmt.Println(string(clairServerOut))

		fmt.Fprintf(w,"Deploying Clair with Persistence... Please wait few seconds and go back")
	}	
}


func clairScan_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		deleteClairServiceOut, deleteClairServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete service clair-svc -n security-scanner").CombinedOutput() 
		if deleteClairServiceErr != nil {
			fmt.Println(deleteClairServiceErr)
		}
		fmt.Println(string(deleteClairServiceOut))

		copyClairServiceOut, copyClairServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-job-service.yaml /tmp/security-scanner-clair-job-service.yaml").CombinedOutput() 
		if copyClairServiceErr != nil {
			fmt.Println(copyClairServiceErr)
		}
		fmt.Println(string(copyClairServiceOut))

		clairServiceOut, clairServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-job-service.yaml").CombinedOutput() 
		if clairServiceErr != nil {
			fmt.Println(clairServiceErr)
		}
		fmt.Println(string(clairServiceOut))

		deleteClairJobOut, deleteClairJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete job security-scanner-clair -n security-scanner").CombinedOutput() 
		if deleteClairJobOutErr != nil {
			fmt.Println(deleteClairJobOutErr)
		}
		fmt.Println(string(deleteClairJobOut))

		time.Sleep(40 * time.Second)

		copyClairJobOut, copyClairJobErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-clair-job.yaml /tmp/security-scanner-clair-job.yaml").CombinedOutput() 
		if copyClairJobErr != nil {
			fmt.Println(copyClairJobErr)
		}
		fmt.Println(string(copyClairJobOut))

		registryServer := string(r.URL.Query().Get("registryServer"))

		registryServerOut, registryServerErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~REGISTRYSERVER_PLACEHOLDER~" + registryServer + "~g' /tmp/security-scanner-clair-job.yaml").CombinedOutput() 
		if registryServerErr != nil {
			fmt.Println(registryServerErr)
		}
		fmt.Println(string(registryServerOut))

		registryUsername := string(r.URL.Query().Get("registryUsername"))

		registryUsernameOut, registryUsernameErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~REGISTRYUSERNAME_PLACEHOLDER~" + registryUsername + "~g' /tmp/security-scanner-clair-job.yaml").CombinedOutput() 
		if registryUsernameErr != nil {
			fmt.Println(registryUsernameErr)
		}
		fmt.Println(string(registryUsernameOut))

		registryPassword := string(r.URL.Query().Get("registryPassword"))
	
		registryPasswordOut, registryPasswordErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~REGISTRYPASSWORD_PLACEHOLDER~" + registryPassword + "~g' /tmp/security-scanner-clair-job.yaml").CombinedOutput() 
		if registryPasswordErr != nil {
			fmt.Println(registryPasswordErr)
		}
		fmt.Println(string(registryPasswordOut))

		registryImage := string(r.URL.Query().Get("registryImage"))
	
		registryImageOut, registryImageErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~REGISTRYIMAGE_PLACEHOLDER~" + registryImage + "~g' /tmp/security-scanner-clair-job.yaml").CombinedOutput() 
		if registryImageErr != nil {
			fmt.Println(registryImageErr)
		}
		fmt.Println(string(registryImageOut))

		clairJobOut, clairJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-clair-job.yaml").CombinedOutput() 
		if clairJobOutErr != nil {
			fmt.Println(clairJobOutErr)
		}
		fmt.Println(string(clairJobOut))

		registryImageWithoutSlashes := strings.ReplaceAll(registryImage, "/", "-")
		registryImageWithoutSlashesAndColon := strings.ReplaceAll(registryImageWithoutSlashes, ":", "-")
		reportName := "analysis-" + registryImageWithoutSlashesAndColon 
		
		if(strings.Contains(registryImage, ":")){
			reportName = reportName + ".html"
		} else {
			reportName = reportName + "-latest.html"
			
		}
		fmt.Println("Report Name "+ reportName)

		podName := ""
		timeout := 100
		for {
			if timeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				podNameOut, podNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Running --selector=app=security-scanner-clair -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if podNameErr != nil {
					fmt.Println(podNameErr)
				}
				podName = string(podNameOut[17:len(podNameOut)])
				podName = strings.TrimSuffix(podName, "\n")	
				
				if len(podName) > 0 {
					break
				} else {
					time.Sleep(1 * time.Second) 
					timeout--  
				}                      
			}
		}

		fmt.Println("Pod " + podName)

		permissionClairResultsJobOut, permissionClairResultsJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && mkdir -p /scanresults && chmod -R 777 /scanresults").CombinedOutput() 
		if permissionClairResultsJobOutErr != nil {
			fmt.Println(permissionClairResultsJobOutErr)
		}
		fmt.Println(string(permissionClairResultsJobOut))
		
		time.Sleep(120 * time.Second)

		copyClairResultOut, clairResultErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl cp " + podName + ":/clairctl/reports/html/" + reportName + " /scanresults/" + reportName + " -n security-scanner").CombinedOutput()
		if clairResultErr != nil {
			fmt.Println(clairResultErr)
		}
		fmt.Println(string(copyClairResultOut))

		fmt.Println("kubectl cp " + podName + ":/clairctl/reports/html/" + reportName + " /scanresults/" + reportName + " -n security-scanner")

		fmt.Fprintf(w,`<h1>Security Scanner - Clair</h1>
		<p><button><a href="scanresults/` + reportName + `" target="_blank" style="text-decoration:none">Show Scan Results</a></button></p>
		`)	
	}
}

func csysdigScan_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		deleteCsysdigDepOut, deleteCsysdigDepOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete -f security-scanner-csysdig-deployment.yaml ").CombinedOutput() 
		if deleteCsysdigDepOutErr != nil {
			fmt.Println(deleteCsysdigDepOutErr)
		}
		fmt.Println(string(deleteCsysdigDepOut))

		time.Sleep(30 * time.Second) 

		out, err := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f security-scanner-csysdig-deployment.yaml").CombinedOutput() 
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		
		podName := ""
		timeout := 30
		for {
			if timeout <= 0 {
				break
			} else {
				if len(podName) > 0 {
					break
				} else {
					time.Sleep(10 * time.Second) 
					timeout--  
				} 
		
				podNameOut, podNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Running --selector=app=security-scanner-csysdig -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if podNameErr != nil {
					fmt.Println(podNameErr)
				}
				podName = string(podNameOut[17:len(podNameOut)])
				podName = strings.TrimSuffix(podName, "\n")	                     
			}
		}

		fmt.Println("Pod " + podName)

		fmt.Fprintf(w,`kubectl exec -it ` + podName + ` -n security-scanner sh`)
	}
}

func kubebenchScan_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		deletekubebenchJobOut, deletekubebenchJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete job security-scanner-kube-bench -n security-scanner").CombinedOutput() 
		if deletekubebenchJobOutErr != nil {
			fmt.Println(deletekubebenchJobOutErr)
		}
		fmt.Println(string(deletekubebenchJobOut))

		deletekubebenchResultsJobOut, deletekubebenchResultsJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && rm -rf /scanresults/security-scanner-kube-bench-job-scan-results.txt").CombinedOutput() 
		if deletekubebenchResultsJobOutErr != nil {
			fmt.Println(deletekubebenchResultsJobOutErr)
		}
		fmt.Println(string(deletekubebenchResultsJobOut))

		copykubebenchJobOut, kubebenchJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-kube-bench-job.yaml /tmp/security-scanner-kube-bench-job.yaml").CombinedOutput() 
		if kubebenchJobOutErr != nil {
			fmt.Println(kubebenchJobOutErr)
		}
		fmt.Println(string(copykubebenchJobOut))

		outKubebench, errKubebench := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-kube-bench-job.yaml").CombinedOutput() 
		if errKubebench != nil {
			fmt.Println(errKubebench)
		}
		fmt.Println(string(outKubebench))

		podName := ""
		timeout := 10
		for {
			if timeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				if len(podName) > 0{
					break
				} else {
					time.Sleep(10 * time.Second) 
					timeout--  
				} 

				podNameOut, podNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Running --selector=app=security-scanner-kube-bench -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if podNameErr != nil {
					fmt.Println(podNameErr)
				}
				podName = string(podNameOut[17:len(podNameOut)])
				podName = strings.TrimSuffix(podName, "\n")	                     
			}
		}
		fmt.Println("Pod " + podName)

		succeededPodName := ""
		succeededTimeout := 50
		for {
			if succeededTimeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				if len(succeededPodName) > 0 {
					break
				} else {
					time.Sleep(1 * time.Second) 
					succeededTimeout--  
				} 

				copyKubebenchResultOut, kubebenchJobResultErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl cp " + podName + ":/scanresults/security-scanner-kube-bench-job-scan-results.txt /scanresults/security-scanner-kube-bench-job-scan-results.txt -n security-scanner").CombinedOutput()
				if kubebenchJobResultErr != nil {
					fmt.Println(kubebenchJobResultErr)
				}
				fmt.Println(string(copyKubebenchResultOut))
		
				succeededPodNameOut, succeededPodNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Succeeded --selector=app=security-scanner-kube-bench -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if succeededPodNameErr != nil {
					fmt.Println(succeededPodNameErr)
				}
				succeededPodName = string(succeededPodNameOut[17:len(succeededPodNameOut)])
				succeededPodName = strings.TrimSuffix(succeededPodName, "\n")	                     
			}
		}

		fmt.Println("Succeeded Pod " + succeededPodName)

		fmt.Fprintf(w,`<h1>Security Scanner - Kube-Bench</h1>
		<p><button><a href="scanresults/security-scanner-kube-bench-job-scan-results.txt" target="_blank" style="text-decoration:none">Show Scan Results</a></button></p>
		`)	
	}
}

func kubehunterScan_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		deletekubehunterJobOut, deletekubehunterJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete job security-scanner-kube-hunter -n security-scanner").CombinedOutput() 
		if deletekubehunterJobOutErr != nil {
			fmt.Println(deletekubehunterJobOutErr)
		}
		fmt.Println(string(deletekubehunterJobOut))

		deletekubehunterResultsJobOut, deletekubehunterResultsJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && rm -rf /scanresults/security-scanner-kube-hunter-job-scan-results.txt").CombinedOutput() 
		if deletekubehunterResultsJobOutErr != nil {
			fmt.Println(deletekubehunterResultsJobOutErr)
		}
		fmt.Println(string(deletekubehunterResultsJobOut))

		copykubehunterJobOut, kubehunterJobOutErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-kube-hunter-job.yaml /tmp/security-scanner-kube-hunter-job.yaml").CombinedOutput() 
		if kubehunterJobOutErr != nil {
			fmt.Println(kubehunterJobOutErr)
		}
		fmt.Println(string(copykubehunterJobOut))

		outKubehunter, errKubehunter := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-kube-hunter-job.yaml").CombinedOutput() 
		if errKubehunter != nil {
			fmt.Println(errKubehunter)
		}
		fmt.Println(string(outKubehunter))

		podName := ""
		timeout := 10
		for {
			if timeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				if len(podName) > 0{
					break
				} else {
					time.Sleep(10 * time.Second) 
					timeout--  
				} 

				podNameOut, podNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Running --selector=app=security-scanner-kube-hunter -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if podNameErr != nil {
					fmt.Println(podNameErr)
				}
				podName = string(podNameOut[17:len(podNameOut)])
				podName = strings.TrimSuffix(podName, "\n")	                     
			}
		}
		fmt.Println("Pod " + podName)

		succeededPodName := ""
		succeededTimeout := 100
		for {
			if succeededTimeout <= 0 {
				fmt.Fprintf(w,`Timeout Error`)
				return 
			} else {
				if len(succeededPodName) > 0 {
					break
				} else {
					time.Sleep(1 * time.Second) 
					succeededTimeout--  
				} 

				copyKubehunterResultOut, kubehunterJobResultErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl cp " + podName + ":/scanresults/security-scanner-kube-hunter-job-scan-results.txt /scanresults/security-scanner-kube-hunter-job-scan-results.txt -n security-scanner").CombinedOutput()
				if kubehunterJobResultErr != nil {
					fmt.Println(kubehunterJobResultErr)
				}
				fmt.Println(string(copyKubehunterResultOut))
		
				succeededPodNameOut, succeededPodNameErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods --field-selector status.phase=Succeeded --selector=app=security-scanner-kube-hunter -o=name -n security-scanner | sed 's/^.\\{4\\}//'").CombinedOutput()
				if succeededPodNameErr != nil {
					fmt.Println(succeededPodNameErr)
				}
				succeededPodName = string(succeededPodNameOut[17:len(succeededPodNameOut)])
				succeededPodName = strings.TrimSuffix(succeededPodName, "\n")	                     
			}
		}

		fmt.Println("Succeeded Pod " + succeededPodName)

		fmt.Fprintf(w,`<h1>Security Scanner - Kube-Hunter</h1>
		<p><button><a href="scanresults/security-scanner-kube-hunter-job-scan-results.txt" target="_blank" style="text-decoration:none">Show Scan Results</a></button></p>
		`)	
	}
}

func falcoDeploy_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		deleteFalcosidekickServiceOut, deleteFalcosidekickServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete service falcosidekick-svc -n security-scanner").CombinedOutput() 
		if deleteFalcosidekickServiceErr != nil {
			fmt.Println(deleteFalcosidekickServiceErr)
		}
		fmt.Println(string(deleteFalcosidekickServiceOut))

		copyFalcosidekickServiceOut, copyFalcosidekickServiceErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-falcosidekick-service.yaml /tmp/security-scanner-falcosidekick-service.yaml").CombinedOutput() 
		if copyFalcosidekickServiceErr != nil {
			fmt.Println(copyFalcosidekickServiceErr)
		}
		fmt.Println(string(copyFalcosidekickServiceOut))

		outFalcosidekickService, errFalcosidekickService := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-falcosidekick-service.yaml").CombinedOutput() 
		if errFalcosidekickService != nil {
			fmt.Println(errFalcosidekickService)
		}
		fmt.Println(string(outFalcosidekickService))

		deleteFalcosidekickOut, deleteFalcosidekickErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete deployment security-scanner-falcosidekick -n security-scanner").CombinedOutput() 
		if deleteFalcosidekickErr != nil {
			fmt.Println(deleteFalcosidekickErr)
		}
		fmt.Println(string(deleteFalcosidekickOut))

		copyFalcosidekickOut, copyFalcosidekickErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-falcosidekick-deployment.yaml /tmp/security-scanner-falcosidekick-deployment.yaml").CombinedOutput() 
		if copyFalcosidekickErr != nil {
			fmt.Println(copyFalcosidekickErr)
		}
		fmt.Println(string(copyFalcosidekickOut))

		fmt.Println(r.URL.Query().Get("teamsWebhookURL"))
		fmt.Println("sed -i 's~TEAMS_WEBHOOKURL~" + r.URL.Query().Get("teamsWebhookURL") + "~g' /tmp/security-scanner-falcosidekick-deployment.yaml")

		SetFalcosidekickOut, falcosidekickSetErr := exec.Command("sh","-c","echo \"security-scanner\" && sed -i 's~TEAMS_WEBHOOKURL_PLACEHOLDER~" + r.URL.Query().Get("teamsWebhookURL") + "~g' /tmp/security-scanner-falcosidekick-deployment.yaml").CombinedOutput() 
		if falcosidekickSetErr != nil {
			fmt.Println(falcosidekickSetErr)
		}
		fmt.Println(string(SetFalcosidekickOut))

		outFalcosidekick, errFalcosidekick := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-falcosidekick-deployment.yaml").CombinedOutput() 
		if errFalcosidekick != nil {
			fmt.Println(errFalcosidekick)
		}
		fmt.Println(string(outFalcosidekick))

		time.Sleep(30 * time.Second)

		deleteFalcoOut, deleteFalcoErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete daemonset security-scanner-falco -n security-scanner").CombinedOutput() 
		if deleteFalcoErr != nil {
			fmt.Println(deleteFalcoErr)
		}
		fmt.Println(string(deleteFalcoOut))

		copyFalcoOut, falcoOutErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-falco-daemonset.yaml /tmp/security-scanner-falco-daemonset.yaml").CombinedOutput() 
		if falcoOutErr != nil {
			fmt.Println(falcoOutErr)
		}
		fmt.Println(string(copyFalcoOut))

		outFalco, errFalco := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-falco-daemonset.yaml").CombinedOutput() 
		if errFalco != nil {
			fmt.Println(errFalco)
		}
		fmt.Println(string(outFalco))

		time.Sleep(50 * time.Second)

		fmt.Fprintf(w,`<h1>Security Scanner - Falco</h1>
		<p>Successfully deployed Falco inside the cluster</p>
		`)	
	}
}


func fetchimages_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		namespace := r.URL.Query().Get("namespace")
		fmt.Println(namespace)

		output_command := `"{..image}"`
		fmt.Println(output_command)

		fmt.Println("kubectl get pods -n '" + namespace + "' -o jsonpath=" + output_command + " |tr -s '[[:space:]]' '\\n' |sort |uniq -c | awk '!visited[$0]++' | awk '{print $2}'")

		fetchimagesOut, fetchimagesErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl get pods -n '" + namespace + "' -o jsonpath=" + output_command + " |tr -s '[[:space:]]' '\\n' |sort |uniq -c | sed 's/:latest//g' | awk '!visited[$0]++' | awk '{print $2}'").CombinedOutput() 
		if fetchimagesErr != nil {
			fmt.Println(fetchimagesErr)
		}
		fmt.Println(string(fetchimagesOut[17:len(fetchimagesOut)]))

		output := string(fetchimagesOut[17:len(fetchimagesOut)])

		fmt.Fprintf(w,output)
	}
}	
	
func enablePersistence_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		deleteSecurityScannerOut, deleteSecurityScannerErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl delete deployment security-scanner -n security-scanner").CombinedOutput() 
		if deleteSecurityScannerErr != nil {
			fmt.Println(deleteSecurityScannerErr)
		}
		fmt.Println(string(deleteSecurityScannerOut))

		copySecurityScannerOut, copySecurityScannerErr := exec.Command("sh","-c","echo \"security-scanner\" && cp security-scanner-deployment-pv.yaml /tmp/security-scanner-deployment-pv.yaml").CombinedOutput() 
		if copySecurityScannerErr != nil {
			fmt.Println(copySecurityScannerErr)
		}
		fmt.Println(string(copySecurityScannerOut))

		SecurityScannerOut, SecurityScannerErr := exec.Command("sh","-c","echo \"security-scanner\" && kubectl apply -f /tmp/security-scanner-deployment-pv.yaml").CombinedOutput() 
		if SecurityScannerErr != nil {
			fmt.Println(SecurityScannerErr)
		}
		fmt.Println(string(SecurityScannerOut))

		fmt.Fprintf(w,"Enabling Persistence... Please wait few seconds and go back")
	}
}
