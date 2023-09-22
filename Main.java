import src.CacheController;

class Main {
    public static void main (String[] args) {
        CacheController cc = new CacheController();
        cc.write(0, 2, 3, "E");
        System.out.println(cc.read(0)[0]);
        System.out.println(cc.read(0)[1]);
        System.out.println(cc.read(0)[2]);
    }

}

