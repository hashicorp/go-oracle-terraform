package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccSSHKeyLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-key"

	sshKeyClient, err := getSSHKeysClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained SSH Key Client\n")

	createSSHKeyInput := CreateSSHKeyInput{
		Name:    name,
		Key:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC7BzZyp8CWN7tfIZiZwWx8H9RO2ClKu0ru/+bGEsUmHmSS7aA+iPBVqeK1Kr2nVkoG/32GaMLfVNRlRgZZGTBTFM5nnytNoo3DC9dnIPuIu95TbF1afGkVFNNyeJkC5bQDDaRDMaYBbPVJTa6bA8v7nmzvElQHPRtdRqZnFx80QHdrgTluqhtrxWDBCYMSm2meL/NU11kijoKfYSReT4lroglSxnkvP0vjUqUSvZ6tI231Ggvxg4TU1TL4OgtNyfQgXK585V05n7IT9iiJHThah2/ZGsb0DZimj/D5LxngciXVOkOR1sDt8pQb7QCxgoxOO3sa1K3pFi5UAJQ10tSyhu0yn0AnRG13NWK6DlLKhLzZM5jhGJeeYYuwCL5fzJojflouHgebOO62gqNANkUcf7cWUBJRWjSAYuXe/C6rJOriZuUkC87QpffpYd2WaJmqnjAaj7NaqOTzk5ltpS39EjMenyXWWw1MPs7eEB/A/Rfol0cHzGqoXaIZAJVaEpW7ePWEj193CqSc6uh1nwAT15rvh63z2l1iPL0CbuF4GwZWsIZ6roirmwPpKY79kAls69EKsa7bydSQuYpbU5otkT20FIbtHmyFMYpJzYM6sQHoljO2AHWmWChkYtqglbFPrQgwIrsAHbJtmzNcmbXLUm1AY+SjZd1UYqPBjFDb7w== your_email@example.com",
		Enabled: true,
	}
	sshKey, err := sshKeyClient.CreateSSHKey(&createSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created SSH Key: %+v\n", sshKey)

	defer tearDownSSHKey(t, sshKeyClient, name)

	getSSHKeyInput := GetSSHKeyInput{
		Name: name,
	}
	getSSHKeyOutput, err := sshKeyClient.GetSSHKey(&getSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	if createSSHKeyInput.Key != getSSHKeyOutput.Key {
		t.Fatalf("Created and retrived keys don't match %s\n%s\n", createSSHKeyInput.Key, getSSHKeyOutput.Key)
	}
	log.Printf("Successfully retrieved ssh key\n")

	updateSSHKeyInput := UpdateSSHKeyInput{
		Name:    name,
		Key:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC7BzZyp8CWN7tfIZiZwWx8H9RO2ClKu0ru/+bGEsUmHmSS7aA+iPBVqeK1Kr2nVkoG/32GaMLfVNRlRgZZGTBTFM5nnytNoo3DC9dnIPuIu95TbF1afGkVFNNyeJkC5bQDDaRDMaYBbPVJTa6bA8v7nmzvElQHPRtdRqZnFx80QHdrgTluqhtrxWDBCYMSm2meL/NU11kijoKfYSReT4lroglSxnkvP0vjUqUSvZ6tI231Ggvxg4TU1TL4OgtNyfQgXK585V05n7IT9iiJHThah2/ZGsb0DZimj/D5LxngciXVOkOR1sDt8pQb7QCxgoxOO3sa1K3pFi5UAJQ10tSyhu0yn0AnRG13NWK6DlLKhLzZM5jhGJeeYYuwCL5fzJojflouHgebOO62gqNANkUcf7cWUBJRWjSAYuXe/C6rJOriZuUkC87QpffpYd2WaJmqnjAaj7NaqOTzk5ltpS39EjMenyXWWw1MPs7eEB/A/Rfol0cHzGqoXaIZAJVaEpW7ePWEj193CqSc6uh1nwAT15rvh63z2l1iPL0CbuF4GwZWsIZ6roirmwPpKY79kAls69EKsa7bydSQuYpbU5otkT20FIbtHmyFMYpJzYM6sQHoljO2AHWmWChkYtqglbFPrQgwIrsAHbJtmzNcmbXLUm1AY+SjZd1UYqPBjFDb7w== your_email@example.com",
		Enabled: false,
	}
	updateSSHKeyOutput, err := sshKeyClient.UpdateSSHKey(&updateSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	if updateSSHKeyOutput.Enabled != updateSSHKeyInput.Enabled {
		t.Fatalf("Key not successfully updated \nDesired: %s \nActual: %s", updateSSHKeyInput.Key, updateSSHKeyOutput.Key)
	}
	log.Printf("Successfully updated ssh key\n")
}

func getSSHKeysClient() (*SSHKeysClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SSHKeysClient{}, err
	}

	return client.SSHKeys(), nil
}

func tearDownSSHKey(t *testing.T, client *SSHKeysClient, name string) {
	deleteSSHKeyInput := DeleteSSHKeyInput{
		Name: name,
	}
	err := client.DeleteSSHKey(&deleteSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted SSH Key")
}
