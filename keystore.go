package main

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
)

/*
keystore 是一个包含经过加密了的钱包私钥。go-ethereum 中的 keystore，
每个文件只能包含一个钱包密钥对。要生成 keystore，首先您必须调用 NewKeyStore，
给它提供保存 keystore 的目录路径。然后，您可调用 NewAccount 方法创建新的钱包，
并给它传入一个用于加密的口令。您每次调用 NewAccount，它将在磁盘上生成新的 keystore 文件。
*/
func createKs() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address.Hex())
}

/*
*
现在要导入您的 keystore，您基本上像往常一样再次调用 NewKeyStore，然后调用 Import 方法，
该方法接收 keystore 的 JSON 数据作为字节。第二个参数是用于加密私钥的口令。第三个参数是指定一个新的加密口令，
但我们在示例中使用一样的口令。导入账户将允许您按期访问该账户，但它将生成新 keystore 文件！
有两个相同的事物是没有意义的，所以我们将删除旧的
*/
func importKs() {
	file := "./tmp/UTC--2024-07-14T14-47-08.572596000Z--9ef729b4419fba4c60c6a296b3581dad0c403530"
	ks := keystore.NewKeyStore("./temp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := "secret"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	b := !errors.Is(os.Remove(file), err)
	if b && err != nil {
		log.Fatal(err)
	}
}

func main() {
	//createKs()
	importKs()
}
