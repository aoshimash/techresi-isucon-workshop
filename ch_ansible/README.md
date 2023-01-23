# Ansible study

## Contents:
1. What is Ansible?
2. Inventory
3. ansible.cfg
4. Playbook
5. Variable priority
6. Ansible vault
7. jinja2 template
8. Ansible role
9. Tips

<br>
  
## 1. What is Ansible?

- Based on the ”Infrastructure as Code” philosophy, we can manage our hosts from ansible server by text file (.yml file).
- We can define hosts managed by ansible server in inventory, and there is no need to install new packages on them since it is agentless.
- Since Ansible is idempotent, we will declare "How it should be". Therefore, the same result will be obtained no matter how many times it is executed.

<br>
  
## 2. Inventory

We can define hosts managed by ansible server in inventory file. 
Here is an example of inventory file.
```
$ cat inventory.ini
  dev1.example.com
  dev2.example.com 
  pro1.example.com 
  pro2.example.com 
```

We can also manage hosts as a inventory group.
```
$ cat inventory.ini
  [webservers]         <<< group name
  dev1.example.com
  pro1.example.com

  [dbservers]
  dev2.example.com
  pro2.example.com

  [development]
  dev1.example.com
  dev2.example.com 

  [production]  
  pro1.example.com 
  pro2.example.com 
```

We can also make group of groups.
```
$ cat inventory.ini
  [development]
  dev1.example.com
  dev2.example.com 

  [production]  
  pro1.example.com 
  pro2.example.com 

  [all:children]     <<< group "all" includes group "development" and grouop "production".
  development 
  production
```

<br>
  
## 3. ansible.cfg

Certain settings in Ansible are adjustable via ansible.cfg.
https://docs.ansible.com/ansible/2.9_ja/reference_appendices/config.html#ansible-configuration-settings  
Here is an sample ansible.cfg file.
```
$ cat ancible.cfg
  [defaults]
  host_key_checking = False
  vault_password_file = /etc/ansible/.vault_pass
  log_path=/var/log/ansible.log
  inventory=inventory.ini
  ansible_roles=~/.ansible/roles:/usr/share/ansible/roles:/etc/ansible/roles

  [privilege_escalation]
  become = true
  become_method = sudo
  become_user = root
  become_ask_pass = False

  [persistent_connection]
  command_timeout = 30
```

Ansible searches for ansible.cfg in these locations in order:

1. use $ANSIBLE_CONFIG (environment variable if set)
2. ansible.cfg (in the current directory)
3. ~/.ansible.cfg (in the home directory)
4. /etc/ansible/ansible.cfg

We can check which ansible.cfg is using by the command "ansible --version"

```
$ ansible --version
  ansible [core 2.12.1]
  config file = /etc/ansible/ansible.cfg    <<<<
  ...
```

<br>
  
## 4. Playbook

Here is an example of playboook.
We can find all modules (ex. "yum", "service" as below) in Ansible Documentation: https://docs.ansible.com/ansible/2.9/modules/list_of_all_modules.html
```
$ cat example.yml
- hosts: webservers                           <<< Target host (We can specify where to install)
  vars:                                       <<< Vars (We can define variables if nessesary)
    - package : httpd
    - service : httpd
  tasks:                                      <<< Tasks (We can run modules)
  - name: Install webpackage
    yum:
      name: "{{package}}"
      state: installed
  - name: Start and enable web service
    service:
      name: "{{service}}"
      state: started
      enabled : yes

- hosts: dbservers
  vars:
  ...
```

When we run the playbook, each name key will be displayed (ex. "Install webpackage", "Start and enable web service" as below)
with the result: ok /changed /unreachable /failed /skipped /rescued /ignored.
Keep in mind it's better to use understandable name.
```
$ ansible-playbook example.yml
PLAY [webservers] **********************************************************************************************
TASK [Gathering Facts] *****************************************************************************************
ok: [X.X.X.X]
TASK [Install webpackage] **************************************************************************************
ok: [X.X.X.X]
TASK [Start and enable web service] ****************************************************************************
changed: [X.X.X.X]
PLAY [dbservers] ***********************************************************************************************
TASK [Gathering Facts] *****************************************************************************************
ok: [X.X.X.X]
...
PLAY RECAP *****************************************************************************************************
localhost : ok=6 changed=2 unreachable=0 failed=0 skipped=0 rescued=0 ignored=0
```

※What is the TASK [Gathering Facts]??  
Ansible facts is the variables automatically discovered by Ansible on the management host.
(ex. hostname, os version, interfaces, memory/disks...)

we can confirm facts by command "ansible [hostname] –m setup"
```
$ ansible all –m setup
```

<br>
  
## 5. Variable priority
Ansible allows us to write many variables in different places, we should understand their priority.
(https://docs.ansible.com/ansible/2.9/user_guide/playbooks_variables.html?highlight=variable%20precedence#variable-precedence-where-should-i-put-a-variable)
Qiita(https://qiita.com/answer_d/items/b8a87aff8762527fb319)

1. command line values (eg “-u user”)
2. role defaults 
3. inventory file or script group vars  
4. inventory group_vars/all 
5. playbook group_vars/all 
6. inventory group_vars/* 
7. playbook group_vars/* 
8. inventory file or script host vars 
9. inventory host_vars/* 
10. playbook host_vars/* 
11. host facts / cached set_facts 
12. play vars
13. play vars_prompt
14. play vars_files
15. role vars (defined in role/vars/main.yml)
16. block vars (only for tasks in block)
17. task vars (only for the task)
18. include_vars
19. set_facts / registered vars
20. role (and include_role) params
21. include params
22. extra vars (always win precedence)

<br>

## 6. Ansible Vault

We can use Ansible Vault to protect sensitive data by encrypting.
<br><br>
- Encrypt the file using Ansible Vault.

```
$ cat secret_file
This file is secret.

$ ansible-vault encrypt secret_file 
New Vault password: Password123!
Confirm New Vault password: Password123!
Encryption successful

$ cat secret_file
$ANSIBLE_VAULT;1.1;AES256
61336364633362643764623434363963633764656331646365653333646639353234313231663163
6133313563343334326338303537306466383939366263640a323436313433386465306635313035
66303233363939666563653236656636396161346161613138326434306337303238663565613033
3935333138623263310a353461323432623263336464346430623238383863653836643862626331
33333030356535666635373931376633633337653335386263343461613734333138
```
<br>

- Decrypt the file using Ansible Vault.
```
 $ cat secret_file
$ANSIBLE_VAULT;1.1;AES256
61336364633362643764623434363963633764656331646365653333646639353234313231663163
6133313563343334326338303537306466383939366263640a323436313433386465306635313035
66303233363939666563653236656636396161346161613138326434306337303238663565613033
3935333138623263310a353461323432623263336464346430623238383863653836643862626331
33333030356535666635373931376633633337653335386263343461613734333138

$ ansible-vault decrypt secret_file
Vault password: Password123!
Decryption successful

$ cat secret_file
This file is secret.
```
<br>

- Encrypt the string using Ansible Vault.
```
$ ansible-vault encrypt_string This_is_secret_password
New Vault password: Password123!
Confirm New Vault password: Password123!
!vault |
          $ANSIBLE_VAULT;1.1;AES256
          36633431666139306637623936643163393065623730373939353865313730316265336332363031
          3633383539383731343235383562666564323034646232610a646134666465343630356538376534
          36343537643532356330363138323166313764626236653935363432353936646263666637636266
          3735383031343338660a306663353737643830623831303939306263396365396434323338343163
          63303830393865646163613339353365343331323031616137363333363666663161
Encryption successful
```
<br>

- After we added the encrypted string to a var file (vars.yml), we can see the original value using the debug module.
```
$ cat vars.yml
new_password: !vault |
          $ANSIBLE_VAULT;1.1;AES256
          36633431666139306637623936643163393065623730373939353865313730316265336332363031
          3633383539383731343235383562666564323034646232610a646134666465343630356538376534
          36343537643532356330363138323166313764626236653935363432353936646263666637636266
          3735383031343338660a306663353737643830623831303939306263396365396434323338343163
          63303830393865646163613339353365343331323031616137363333363666663161

$ ansible localhost -m debug -a var="new_password" -e "@vars.yml" --ask-vault-pass
Vault password: Password123!
localhost | SUCCESS => {
    "new_password": "This_is_secret_password"
```

<br>

## 7. Jinja2 template

We can copy a jinja2 format file with variables embedded with the template module.<br>
Here is an example of how to use jinja2 template file.

playbook
```
...
- name: jinja2 template test
  hosts:
    - host1.examle.com
    - host2.examle.com
  tasks:
    - name: copy the template file
      template:
        src: files/example.j2     <<< jinja2 template file 
        dest: /etc/example 
        owner: root
        group: root
        mode: 0644
        ...
```

jinja2 template file
```
$ cat files/example.j2
Greetings from {{ inventory_hostname }}     <<< vars(in this case, using Facts vars)
```

result
```
$ ssh host1.example.com
$ cat /etc/example
Greetings from host1.examle.com
```
```
$ ssh host2.examle.com
$ cat /etc/example
Greetings from host2.examle.com
```

<br>

## 8. Ansible Role

If we have a lot of playbooks and want to organize them, we can use Ansible Role to divide the playbooks into those roles and combine / reuse many times for our purposes.

Here is example of project structure:
```
site.yml
site2.yml
site3.yml
roles/
    sample_role_1/
        tasks/
        handlers/
        files/
        templates/
        defaults/
        vars/
        meta/
    sample_role_2/
        tasks/
        handlers/
        files/
        templates/
        defaults/
        vars/
        meta/
        ...
```

- tasks - contains the main list of tasks to be executed by the role.
- handlers - contains handlers, which may be used by this role or even anywhere outside this role.
- files - contains files which can be deployed via this role.
- templates - contains templates which can be deployed via this role.
- defaults - default variables for the role.
- vars - other variables for the role.
- meta - defines some meta data for this role. 

<br>

For example, let's change this tasks to Ansible role.
```
- name: jinja2 template test
  hosts:
    - host1.examle.com
    - host2.examle.com
  tasks:
    - name: copy the template file
      template:
        src: files/example.j2 
        dest: /etc/example 
        owner: root
        group: root
        mode: 0644
        ...
```

↓ Change to Ansible role...

- playbook
```
$cat site.yml
- hosts:
    - host1.examle.com
    - host2.examle.com
  roles:
    - sample_role_1
    - sample_role_2
    ...
```
- sample_role_1/task/main.yml
```
$cat roles/sample_role_1/task/main.yml
- name: copy the template file
  template:
    src: example.j2 
    dest: /etc/example 
    owner: root
    group: root
    mode: 0644
```
- sample_role_1/template/example.j2
```
$cat roles/sample_role_1/template/example.j2
Greetings from {{ inventory_hostname }}
```

<br>

### ※ ansible galaxy

Ansible Galaxy (https://galaxy.ansible.com/) is a public library of Ansible role content. We can search for typical roles and download them locally to use for our purposes.

<br>

## 9. Tips

1. Add the following line to indent two spaces when the Tab key is pressed.
```
$vim ~/.vimrc
autocmd FileType yaml setlocal ai ts=2 sw=2 et
```
<br>

2. We can use "loop" to avoid writing multiple tasks that use the same module.

```
- name: Create user taro
  user:
    name: taro
    state: present
- name: Create user jiro
  user:
    name: jiro
    state: present
- name: Create user saburo
  user:
    name: saburo
    state: present
```

↓ For example, simplify this tasks by using loop...

```
- name: Create users 
  user:
    name: "{{ item }}"
    state: present
  loop:
    - taro
    - jiro
    - saburo
    ...
```
<br>

3. We can use "when" to run tasks in certain conditions.
```
- name: Install packages when freemem is over 100m
    yum:
      name: httpd
      state: installed
    when: ansible_memfree_mb > 100
```
<br>

4. We can use "handlers" to run tasks when only result is "changed".
```
tasks:
  - name: copy example.conf configuration template    
    template:
      src: templates/example.conf
      dest: /etc/httpd/conf.d/example.conf
    notify:
      - restart apache
    handlers:
      - name: restart apache
        service:
          name: httpd 
          state: restarted
```
<br>

5. We can use "register" to temporarily store the result in a variable. We can use "debug" to display the message when run the playbook.

```
$ cat dedug.yml
- hosts: 127.0.0.1
  tasks:
    - name: exec whoami
      shell: whoami
      register: result      <<< store the result in a variable.
    - name: debug result var
      debug:
        msg: "{{result}}"    <<< debug message when run the playbook.


 $ ansible-playbook dedug.yml
PLAY [127.0.0.1] *******************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************
ok: [127.0.0.1]

TASK [exec whoami] *****************************************************************************************************
changed: [127.0.0.1]

TASK [debug result var] ************************************************************************************************
ok: [127.0.0.1] => {
    "msg": {                    <<< debug: {{result}} is shown.
        "changed": true,
        "cmd": "whoami",
        "delta": "0:00:00.011802",
        "end": "2022-02-02 14:02:58.433723",
        "failed": false,
        "msg": "",
        "rc": 0,
        "start": "2022-02-02 14:02:58.421921",
        "stderr": "",
        "stderr_lines": [],
        "stdout": "haruka.takigawa",
        "stdout_lines": [
            "haruka.takigawa"
        ]
    }
}

PLAY RECAP *************************************************************************************************************
127.0.0.1                  : ok=3    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

```

6. DRY RUN
https://docs.ansible.com/ansible/2.9_ja/user_guide/playbooks_checkmode.html