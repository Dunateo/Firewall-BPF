import tkinter as tk
import os
import socket
from threading import Thread
from time import sleep,time

class MaListBox(tk.Tk):

    canvas_id = 0
    numP = 0
    numA = 0
    numB = 0
    id_log = 0

    def __init__(self):
        # Instanciation fenêtre Tk.
        tk.Tk.__init__(self)
        self.size = 750
        self.creer_widgets()
        self.running = 0
        self.addr = None
        self.conn = None
        

    def creer_widgets(self):

        # Lancement du script shell
        # os.system("/home/NameProc.sh") 

        # Ecoute sur le port 234
        # self.ListeningSocket()

        # création canevas Principal
        self.canv = tk.Canvas(self, bg="light gray", height=self.size,
                              width=self.size)
        self.canv.pack(side=tk.LEFT)
        
        global canvas_id
        canvas_id = self.canv.create_text(10, 10, anchor="nw")
        
        # création canevas Log
        self.canv.canvLog = tk.Canvas(self, bg="light blue", height=self.size/3,
                              width=self.size)
        self.canv.canvLog.pack(side=tk.BOTTOM)
        
        global id_log
        id_log = self.canv.canvLog.create_text(10, 10, anchor ="nw")

        # self.openFile()

        # self.dispListProc()

        # self.dispListApp()

        # self.dispBlacklist()

        # Bouton Quitter
        self.bouton_quitter = tk.Button(self, text="Quitter", command=self.quit)
        self.bouton_quitter.pack(side=tk.BOTTOM)
        
        # Bouton Refresh
        # self.bouton_refresh = tk.Button(self, text="Rafraîchir", command=self.refresh)
        # self.bouton_refresh.pack(side=tk.BOTTOM)

        # Bouton Activer / Désactiver Logs
        self.bouton_stop = tk.Button(self, text="Désactiver Logs", command=self.stopc)
        self.bouton_stop.pack(side=tk.BOTTOM)
        self.bouton_start = tk.Button(self, text="Activer Logs", command=self.startc)
        self.bouton_start.pack(side=tk.BOTTOM)
    

    def refresh(self):
        # Lancement du script shell
        os.system("/home/NameProc.sh")
        self.listbox.destroy()
        self.listboxApp.destroy()
        self.listboxBlacklist.destroy()
        self.dispListProc()
        self.dispListApp()
        self.dispBlacklist()
    
    def dispListProc(self):
        global numP
        #self.listbox.destroy()
        # Ouvrir le fichier en lecture seule
        filename = "/home/app.txt"
        file = open(filename, "r")
        # utiliser readlines pour lire toutes les lignes du fichier
        # La variable "lignes" est une liste contenant toutes les lignes du fichier
        lines = file.readlines()
        # On compte le nombre de lignes
        numP = len(lines)
        # Fermez le fichier après avoir lu les lignes
        file.close()
        # Création de la liste
        self.listbox = tk.Listbox(self, height=numP, width=30)
        self.listbox.pack()
        # Itérer sur les lignes
        for line in lines:
            self.listbox.insert(tk.END, line)
        # Selectionne premier élément de listbox.
        self.listbox.select_set(0)
        # Lier une méthode quand clic sur listbox.
        self.listbox.bind("<<ListboxSelect>>", self.clic_listbox)

    def dispListApp(self):
        # Ouvrir le fichier en lecture seule
        filename = "/home/prog.txt"
        file = open(filename, "r")
        # utiliser readlines pour lire toutes les lignes du fichier
        # La variable "lignes" est une liste contenant toutes les lignes du fichier
        lines = file.readlines()
        # On compte le nombre de lignes
        global numA
        numA = len(lines)
        # Fermez le fichier après avoir lu les lignes
        file.close()
        # Création de la liste
        self.listboxApp = tk.Listbox(self, height=numA, width=30)
        self.listboxApp.pack()
        # Itérer sur les lignes
        for line in lines:
            self.listboxApp.insert(tk.END, line)
        # Lier une méthode quand clic sur listbox.
        self.listboxApp.bind("<<ListboxSelect>>", self.clic_listbox)

    def dispBlacklist(self):
        # Ouvrir le fichier en lecture seule
        filename = "/home/oem/Documents/BPF/FromFile/Firewall-BPF-Blocked_ports_from_file/blockedProcess"
        file = open(filename, "r")
        # utiliser readlines pour lire toutes les lignes du fichier
        # La variable "lignes" est une liste contenant toutes les lignes du fichier
        lines = file.readlines()
        # On compte le nombre de lignes
        global numB
        numB = len(lines)
        # Fermez le fichier après avoir lu les lignes
        file.close()
        # Création de la liste
        self.listboxBlacklist = tk.Listbox(self, height=numB, width=30)
        self.listboxBlacklist.pack()
        # Itérer sur les lignes
        for line in lines:
            self.listboxBlacklist.insert(tk.END, line)
        # Lier une méthode quand clic sur listbox.
        self.listboxBlacklist.bind("<<ListboxSelect>>", self.clic_listboxDelete)
        

    def openFile(self):
        filename = "/home/oem/Documents/BPF/FromFile/Firewall-BPF-Blocked_ports_from_file/blockedProcess"
        fichier = open(filename, "r")
        content = fichier.read()
        fichier.close()
        self.canv.itemconfig(canvas_id, text=content)

    def writeFile(self, choix_select):
        # Ecriture dans un fichier
        filename = "/home/oem/Documents/BPF/FromFile/Firewall-BPF-Blocked_ports_from_file/blockedProcess"
        fichier = open(filename, "r+")
        content = fichier.read()
        # Vérification présence choix_select
        if content.find(str(choix_select)) == -1:
            fichier.write(str(choix_select))
            fichier.close()
            self.openFile()
            self.listboxBlacklist.destroy()
            self.dispBlacklist()
        else:
            print("Déjà présent dans la blacklist")
    
    def clic_listbox(self, event):
        # Récup du widget à partir de l'objet event.
        widget = event.widget
        # Récup du choix sélectionné dans la listbox (tuple).
        # (par exemple renvoie `(5,)` si on a cliqué sur `5`)
        selection = widget.curselection()
        # Récup du nombre sélectionné (déjà un entier).
        choix_select = widget.get(selection[0])
        # Affichage sélection
        print("Le choix sélectionné est {}, son type est {}"
              .format(choix_select, type(choix_select)))
        # Gestion écriture fichier
        self.writeFile(choix_select)

    def clic_listboxDelete(self, event):
        # Récup du widget à partir de l'objet event.
        widget = event.widget
        # Récup du choix sélectionné dans la listbox (tuple).
        # (par exemple renvoie `(5,)` si on a cliqué sur `5`)
        selection = widget.curselection()
        # Récup du nombre sélectionné (déjà un entier).
        choix_select = widget.get(selection[0])
        # Gestion écriture fichier
        self.DeleteFile(choix_select)
        self.openFile()
        self.listboxBlacklist.destroy()
        self.dispBlacklist()

    def DeleteFile(self, choix_select):
        filename = "/home/oem/Documents/BPF/FromFile/Firewall-BPF-Blocked_ports_from_file/blockedProcess"
        infile= open(filename,'r')
        lines= infile.readlines() #converting all lines to listelements
        infile.close()
        # make new list, consisting of the words except the one to be removed
        newlist=[i for i in lines if i!=str(choix_select)]  #list comprehension 
        outfile= open(filename,'w')
        outfile.write("".join(newlist))
        outfile.close
        
    def ListeningSocket(self):
        Sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        Host = '127.0.0.1' # l'ip locale de l'ordinateur
        Port = 234         # choix d'un port
        
        # on bind notre socket :
        Sock.bind((Host,Port))
        
        # On est a l'ecoute d'une seule et unique connexion :
        Sock.listen(1)
        
        # Le script se stoppe ici jusqu'a ce qu'il y ait connexion :
        client, adresse = Sock.accept() # accepte les connexions de l'exterieur
        print("L'adresse ",adresse," vient de se connecter au serveur !")
        while 1:
            RequeteDuClient = client.recv(255) # on recoit 255 caracteres grand max
            if not RequeteDuClient: # si on ne recoit plus rien
                    break  # on break la boucle (sinon les bips vont se repeter)
            print(RequeteDuClient,"\a")         # affiche les donnees envoyees, suivi d'un bip sonore
    
    def socket_thread(self):
        print("Thread started..")
        self.ls = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        port = 9999
        self.ls.bind(('', port))
        print("Server listening on port %s" %port)
        self.ls.listen(1)
        self.ls.settimeout(5)
        self.conn=None
        while self.running != 0:
            if self.conn is None:
                try:
                    (self.conn, self.addr) = self.ls.accept()
                    print("Client is at", self.addr[0], "on port", self.addr[1])
                    self.connectionl.configure(text="CONNECTED!")

                except socket.timeout as e:
                    print ("Waiting for Connection...")

                except Exception as e:
                    print("Connect exception: "+str(e) )

            if self.conn != None:
                print ("Connected to "+str(self.conn)+","+str(self.addr))
                self.conn.settimeout(5)
                self.rc = ""
                connect_start = time() # actually, I use this for a timeout timer
                while self.rc != "done":
                    self.rc=''
                    try:
                        content = self.canv.canvLog.itemcget(id_log, 'text')
                        if len(content):
                            self.rc = self.conn.recv(1000).decode('utf-8')
                            total = content + '\n' + self.rc
                            self.canv.canvLog.itemconfig(id_log, text=total)
                        else:
                            self.rc = self.conn.recv(1000).decode('utf-8')
                            self.canv.canvLog.itemconfig(id_log, text=self.rc)

                    except Exception as e:
                        # we can wait on the line if desired
                        print ("socket error: "+repr(e))

                    if len(self.rc):
                        print("Got data :", self.rc)
                        msg = "Got data.\n"
                        strTobyte = msg.encode('utf-8')
                        self.conn.send(strTobyte)
                        connect_start=time()  # reset timeout time
                    elif (self.running==0) or (time()-connect_start > 30):
                        print ("Tired of waiting on connection!")
                        self.rc = "done"

                print ("Closing connection")
                self.connectionl.configure(text="Not connected.")
                self.conn.close()
                self.conn=None
                print ("connection closed.")

        print ("closing listener...")
        # self running became 0
        self.ls.close()
    
    def startc(self):
        if self.running == 0:
            print ("Starting thread")
            self.running = 1
            self.thread=Thread(target=self.socket_thread)
            self.thread.start()
        else:
            print ("Thread already started.")

    def stopc(self):
        if self.running:
            print("Stopping thread...")
            self.running = 0
            print("TEST milieu...")
            self.thread.join()
            print("TEST fin...")
        else:
            print("Thread not running")


if __name__ == "__main__":
    app = MaListBox()
    app.title("NetCop - GUI")
    app.mainloop()
