package src;

public class Cache {

    private Integer[] data = {0,0,0,0};
    private Integer[] address = {0,0,0,0};
    private String[] state = {"","","",""};

    public void setData(Integer pos, Integer res){
        data[pos] = res;
    }

    public void setAddress(Integer pos, Integer res){
        address[pos] = res;
    }

    public void setState(Integer pos, String res){
        state[pos] = res;
    }

    public Integer getData(Integer pos){
        return data[pos];
    }

    public Integer getAddress(Integer pos){
        return address[pos];
    }

    public String getState(Integer pos){
        return state[pos];
    }
}
