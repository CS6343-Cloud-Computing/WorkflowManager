lines = []
with open("out.txt") as file_in:
    for line in file_in:
        line = line.split(".For")[0].split("id:")[1].strip()
        lines.append(line)
        print(line)
f = open("out1.txt", "w")
for line in lines:
    f.write(line+"\n")
f.close()
