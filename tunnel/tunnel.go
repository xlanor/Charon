package tunnel

import(
	charonConf "xlanor/charon/config"
	logger "xlanor/charon/logger"
	errors
	"golang.org/x/crypto/ssh"
)

func parsePrivateKey() (ssh.Signer, error) {
	buff := charonConf.GetPrivateKey()
	return ssh.parsePrivateKey(buff)
}

func getSshConfig()(*ssh.ClientConfig, error) {
	key, err := parsePrivateKey()
	if err != nil {
		logger.Sugar().Error("Unable to parse private key file")
		return nil, errors.New("Unable to parse private key file")
	}
	cfg := ssh.ClientConfig(
		User: charonConf.GetJumpHostUser(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		}
	)

	return &config, nil
}

func GetPort() *int {
	validate := func(input string) error {
		num, err := strconv.ParseInt(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		if num < 1 || num > 65535 {
			return errors.New("Invalid Port")
		}
		return nil
	}
	while (1){
		prompt := promptui.Prompt{
			Label:    "Enter Local Port",
			Validate: validate,
		}
	
		result, err := prompt.Run()
		if err != nil {
			logger.Sugar().Error("Prompt failed")
			return nil
		}else{
			port, err := CheckPort(result)
			if err != nil {
				logger.Sugar().Error(err.Error())
			}
			if port == nil {
				logger.Sugar().Error("Recevied null pointer in port")
			}else{
				return port
			}
		}
	}

}

func CheckPort(port int) (*int, error){
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// should not ever trigger because of validate
		return nil, err
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer listen.Close()
	return &l.Addr().(*net.TCPAddr).Port, nil
}