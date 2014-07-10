package main
import "crypto/hmac"
import "crypto/sha1"
import "encoding/base64"
import "io/ioutil"
import "mime/multipart"
import "bytes"
import "net/http"
import "fmt"
import "os"

func s3_test() {
    if len(os.Args) < 2 {
        fmt.Println("Enter package name to upload");
        return;
    }

    var client http.Client;

    policy, err := ioutil.ReadFile("./s3.policy");
    if err != nil {
        fmt.Println("Cannot locate s3_policy.txt, required for upload");
        return;
    }

    key, err2 :=  ioutil.ReadFile("./s3.key");
    if err2 != nil {
        fmt.Println("Cannot locate s3.key, required for upload");
        return;
    }

    fmt.Println("* Packing current directory");
    tar_content, err3 := PackCurrentWorkingDir();
    if err3 != nil {
        fmt.Println("Error Unable to pack current Directory");
        fmt.Println(err3);
        return;
    }
    fmt.Println("* Completed packing current directory");

    policy_base64 := base64.StdEncoding.EncodeToString(policy);
    mac := hmac.New(sha1.New, key[:len(key)-1]);
    mac.Write([]byte(policy_base64))
    sum := mac.Sum(nil);
    signature := base64.StdEncoding.EncodeToString(sum[:]);

    buff := new(bytes.Buffer)
    form := multipart.NewWriter(buff);
    form.WriteField("key", "${filename}");
    form.WriteField("AWSAccessKeyId", "AKIAIRCSOT35C7VGHHDA");
    form.WriteField("acl", "public-read");
    form.WriteField("success_action_status", "201");
    form.WriteField("policy", policy_base64);
    form.WriteField("signature", signature);
    form.WriteField("Content-Type", "application/octet-stream");
    file_part, _ := form.CreateFormFile("file", os.Args[1]);
    file_part.Write(tar_content);
    form.WriteField("submit", "Upload to Amazone S3");
    form.Close();

    req, err4 := http.NewRequest("POST","http://libdevdev.s3.amazonaws.com/", buff);
    if err4 != nil {
        fmt.Println("unable to build HTTP request, recived error");
        fmt.Println(err4);
        return;
    }
    req.Header["Content-Type"] = []string{"multipart/form-data; boundary="+form.Boundary()};
    fmt.Println("* Uploading code base to amazone s3");
    response, err5 := client.Do(req);
    if err5 != nil {
        fmt.Println("Unable to get response recived error");
        fmt.Println(err5);
        return;
    }
    answer, err6 := ioutil.ReadAll(response.Body);
    if err6 != nil {
        fmt.Println("Unable to read response body, recived error ");
        fmt.Println(err6);
        return;
    }
    response.Body.Close();
    fmt.Println(string(answer));
}

