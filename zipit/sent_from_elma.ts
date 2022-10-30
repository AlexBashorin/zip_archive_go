
async function toBase64(file: FileItemRef): Promise<unknown> {
    const one_file = await file.fetch()
    const fileObjOne = await fetch(await one_file!.getDownloadUrl());
    const contentOne = new Uint8Array(await fileObjOne.arrayBuffer());
    let binaryOne = '';
    for (const char of contentOne) {
        binaryOne += String.fromCharCode(char);
    }
    const base64One = btoa(binaryOne);
    return base64One
}

async function zipit(): Promise<void> {
    try {
        let firstFileName: unknown = undefined;
        let firstFile: unknown = undefined;
        if(Context.data.tozip_one) {
            const f_name = await Context.data.tozip_one.fetch()
            firstFileName = `${f_name.data.__name}`
            firstFile = await toBase64(Context.data.tozip_one)
        } 

        let secondFileName: unknown = undefined;
        let secondFile: unknown = undefined;
        if(Context.data.tozip_two) {
            const f_name = await Context.data.tozip_two.fetch()
            secondFileName = `${f_name.data.__name}`
            secondFile = await toBase64(Context.data.tozip_two)
        }

        let sent = [
            {
                Name: firstFileName,
                Body: firstFile
            },
            {
                Name: secondFileName,
                Body: secondFile
            },
        ];

        const answer = await fetch("http://192.168.29.174:5050/zipit", {
            method: "POST",
            body: JSON.stringify(sent)
        });

        const g = await answer.text()
        if(g) {
            const text = JSON.stringify(g);

            const binaryAnswer = atob(text);
            let length = binaryAnswer.length;
            const u8arr = new Uint8Array(length);
            while (length--) {
                u8arr[length] = binaryAnswer.charCodeAt(length);
            }
            let fileName = 'zipit.zip'
            Context.data.archive = await Context.fields.archive.create(fileName, u8arr);
        }
    } catch (e) {
        Context.data.error = e;
    }
}