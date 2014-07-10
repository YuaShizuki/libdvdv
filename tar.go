package main
import "archive/tar"
import "compress/gzip"
import "bytes";
import "io/ioutil"

type file_main struct {
    file_name string;
    file_body []byte;
}

func pack_dir(dir string, w *tar.Writer) error {
    file_info, err := ioutil.ReadDir("."+dir);
    if err != nil {
        return err;
    }
    for i := 0; i < len(file_info); i++ {
        if file_info[i].IsDir() {
            if err := pack_dir(dir+"/"+file_info[i].Name(), w); err != nil {
                return err;
            }
            continue;
        }
        file_content,err := ioutil.ReadFile("."+dir+"/"+file_info[i].Name());
        if err != nil {
            return err;
        }
        hdr := &tar.Header{ Name:dir+"/"+file_info[i].Name(),
                            Size:int64(len(file_content)), Mode:int64(file_info[i].Mode()) };
        if err := w.WriteHeader(hdr); err != nil {
            return err;
        }
        if _, err := w.Write(file_content); err != nil {
            return err;
        }
    }
    return nil;
}


func PackCurrentWorkingDir() ([]byte, error) {
    buff := new(bytes.Buffer);
    tw := tar.NewWriter(buff);
    err := pack_dir("", tw);
    if err != nil {
        return nil, err;
    }
    if err := tw.Close(); err != nil {
        return nil, err;
    }
    final_buf := new(bytes.Buffer);
    gw := gzip.NewWriter(final_buf);
    gw.Write(buff.Bytes());
    gw.Close();
    return final_buf.Bytes(), nil;
}
