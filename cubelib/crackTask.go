package cubelib

import (
	"context"
	"cube/model"
	Plugins "cube/plugins"
	"fmt"
	"sync"
	"time"
)

//func GenerateCrackTasks(ip []string, port string, auths []model.Auth, plugins []string) (tasks []model.CrackTask) {
//	tasks = make([]model.CrackTask, 0)
//	for _, i := range ip {
//		for _, auth := range auths {
//			for _, p := range plugins {
//				s := model.CrackTask{Ip: i, Port: port, Auth: &auth, CrackPlugin: p}
//				tasks = append(tasks, s)
//			}
//		}
//	}
//	return tasks
//}

func unitTask(ip string, auths []model.Auth, plugins []string) (tasks []model.CrackTask) {
	tasks = make([]model.CrackTask, 0)
	for _, auth := range auths {
		for _, p := range plugins {
			s := model.CrackTask{Ip: ip, Auth: auth, CrackPlugin: p}
			tasks = append(tasks, s)
		}
	}
	return tasks
}

func processArgs(opt *model.CrackOptions) ([]string, error) {

	return nil, nil
}

func generateAuth(user []string, password []string) (authList []model.Auth) {
	authList = make([]model.Auth, 0)
	for _, u := range user {
		for _, pass := range password {
			a := model.Auth{User: u, Password: pass}
			authList = append(authList, a)
		}
	}
	return authList
}

func saveCrackReport(taskResult model.CrackTaskResult) {
	if len(taskResult.Result) > 0 {
		s := fmt.Sprintf("[*]: %s\n[*]: %s:%d\n", taskResult.CrackTask.CrackPlugin, taskResult.CrackTask.Ip, taskResult.CrackTask.Port)
		s1 := fmt.Sprintf("[*]: %s", taskResult.Result)
		fmt.Println(s + s1)
	}
}

func executeCrackTask(taskChan chan model.CrackTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		//fmt.Println("Hello")
		fn := Plugins.CrackFuncMap[task.CrackPlugin]
		saveCrackReport(fn(task))
	}

}

func RunCrackTasks(tasks []model.CrackTask, scanNum int, timeout int) {
	tasksChan := make(chan model.CrackTask, scanNum*2)
	var wg sync.WaitGroup

	//消费者
	wg.Add(scanNum)
	for i := 0; i < scanNum; i++ {
		go executeCrackTask(tasksChan, &wg)
	}

	//生产者
	//go func() {
	//
	//}()

	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	waitTimeout(&wg, time.Duration(timeout)*time.Second)
}

func StartCrackTask(opt *model.CrackOptions, globalopts *model.GlobalOptions) {

}

func executeUnitTask(taskChan chan model.CrackTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		//fmt.Println("Hello")
		fn := Plugins.CrackFuncMap[task.CrackPlugin]
		saveCrackReport(fn(task))
	}

}

func runCrackTask(ctx context.Context, taskChan chan model.CrackTask) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			fn := Plugins.CrackFuncMap[task.CrackPlugin]
			r := fn(task)
			if len(r.Result) > 0 {
				return
			}
			saveCrackReport(r)
		}
	}
}

func runCrack(plugins []string, ips []string, authList []model.Auth, g model.GlobalOptions) {
	//ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()
	for _, ip := range ips {
		taskChan := make(chan model.CrackTask)

		for i := 0; i < g.Threads; i++ {
			go runCrackTask(ctx, taskChan)
		}

		tasks := unitTask(ip, authList, plugins)
		for _, task := range tasks {
			taskChan <- task
		}

	}
}