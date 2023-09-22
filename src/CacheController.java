package src;

public class CacheController{

    private Cache cache = new Cache(); 

    public String[] read(Integer pos){
        String[] res = new String[3];
        res[0] = cache.getData(pos).toString();
        res[1] = cache.getAddress(pos).toString();
        res[2] = cache.getState(pos);
        return res;
    }

    public void write(Integer pos, Integer data, Integer address, String state){
        if(data >= 0){
            cache.setData(pos, data);
        }

        if(address >= 0){
            cache.setAddress(pos, address);
        }

        if(state != ""){
            cache.setState(pos, state);
        }
    }


}
