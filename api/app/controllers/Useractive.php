<?php

class UseractiveController extends BaseController{
    /**
     * 构造函数
     */
    public $page_num;
    public $coinId;
    public $coinName = '';
    public $activeFee = 0;
    public $activeSend = 0;
    public function init() {
        parent::init();
        if(empty($this->uid)) {
//            $this->echo_json(300,'登录已失效');
        }
        $this->page_num = 10;
        $this->coinId = 6;
        $this->coinName = $this->db->get("member_wallet_coin",'coinname',['id'=>$this->coinId]);
        $this->activeFee = $this->db->get("setting",'content',['id'=>'28']);
        $this->activeSend = $this->db->get("setting",'content',['id'=>'29']);
    }

    //测试时显示手机验证码
//    public function showCodeAction(){
//        $phone = trim($this->requestParams['phone']);
//        $this->echo_json(1,'ok',$this->redis->get("sms_".$phone));
//    }

    //激活账户
    public function activeAccountAction(){
        $num1 = $this->activeFee;
        $num2 = $this->activeSend;
        $username = trim($this->requestParams['username']);
        $password = trim($this->requestParams['password']);
        $activeUsername = trim($this->requestParams['active_username']);
        $deal_pwd = trim($this->requestParams['deal_pwd']);
        $type = intval($this->requestParams['type']);  //激活类型  1 网页激活  2 app激活
        if($type != 1 && $type != 2){
            $this->echo_json(0,'参数错误');
        }
        if($type == 1){
            //网页激活
            if(empty($username) || empty($password) || empty($activeUsername)){
                $this->echo_json(0,'请填写完整信息');
            }
            $user = $this->db->get("member",'*',['username'=>$username]);
            if(empty($user)){
                $this->echo_json(0,'用户不存在');
            }
            if($user['status'] != 1){
                $this->echo_json(0,'该账户无权限激活其他账号');
            }
            //验证密码
            if (!$this->password_verify($password, $user['password'], $user['salt'])) {
                $this->echo_json(0, '密码错误');
            }
            $activeUser = $this->db->get("member",'*',['username'=>$activeUsername]);
            if(empty($activeUser)){
                $this->echo_json(0, '待激活用户信息错误');
            }
            if($activeUser['status'] != 2){
                $this->echo_json(0, '该账户已激活');
            }
        }else{
            //app激活
            if(!$this->uid || !$deal_pwd){
                $this->echo_json(0, '参数错误');
            }
            $user = $this->db->get("member",'*',['uid'=>$this->uid]);
            if(empty($user)){
                $this->echo_json(0,'用户不存在');
            }
            if($user['status'] != 1){
                $this->echo_json(0,'该账户无权限激活其他账号');
            }
            $activeUser = $this->db->get("member",'*',['OR'=>['username'=>$activeUsername,'addr'=>$activeUsername]]);
            if(empty($activeUser)){
                $this->echo_json(0, '待激活用户信息错误');
            }
            if($activeUser['status'] != 2){
                $this->echo_json(0, '该账户已激活');
            }
            $userinfo = $this->db->get('member_wallet_userinfo','*',['uid'=>$this->uid]);
            if(!$this->password_verify($deal_pwd,$userinfo['deal_psw'], $user['salt'])){
                $this->echo_json(0, '支付密码错误');
            }
        }
        $wallet = $this->db->get("member_wallet",'*',['uid'=>$user['uid'],'coin_id'=>$this->coinId]);
        if($num1 > 0){
            if(empty($wallet) || $wallet['total'] < $num1){
                $this->echo_json(0, '余额不足');
            }
        }
        $member = $this->db->get("member",'*',['accountType'=>2]);
        if($member){
            $walletTotal = $this->db->get("member_wallet",'*',['uid'=>$member['uid'],'coin_id'=>$this->coinId]);
        }
        $res = $this->db->action(function() use($num1,$wallet,$activeUser,$num2,$user,$walletTotal){
            try{
                if($num1 > 0){
                    //减金额
                    $res = $this->db->update("member_wallet",['total[-]'=>$num1],['id'=>$wallet['id']]);
                    if(!$res){
                        throw new Exception("减金额失败");
                    }
                    //资金记录
                    $res = $this->db->insert("member_wallet_coin_log",array(
                        'uid'=>$user['uid'],
                        'coin_id'=>$this->coinId,
                        'coin_name'=>$this->coinName,
                        'type'=>102,
                        'note'=>'激活下级',
                        'value'=> -1 * $num1,
                        'balance'=>$wallet['total'],
                        'wallet_id'=>$wallet['id'],
                        'status'=>1,
                        'inputtime'=>time(),
                        'cid'=>$activeUser['uid']
                    ));
                    if(!$res){
                        throw new Exception("加资金记录失败");
                    }
                }
                //修改用户状态
                $res = $this->db->update("member",['status'=>1],['uid'=>$activeUser['uid']]);
                if(!$res){
                    throw new Exception("修改用户状态失败");
                }
                if($num2 > 0 && $walletTotal['total'] >= $num2){

                    //主账号减推广额度
                    $update['total[-]'] = $num2;
                    $res = $this->db->update("member_wallet",$update,['id'=>$walletTotal['id']]);
                    if(!$res){
                        throw new Exception("主账户减金额失败");
                    }
                    $walletTotal['total'] -=  $num2;
                    $res = $this->db->insert("member_wallet_coin_log",array(
                        'uid'=>$walletTotal['uid'],
                        'coin_id'=>$this->coinId,
                        'coin_name'=>$this->coinName,
                        'type'=>103,
                        'note'=>'激活账号赠送',
                        'value'=> -1 * $num2,
                        'balance'=> $walletTotal['total'],
                        'wallet_id'=> $walletTotal['id'],
                        'status'=>1,
                        'inputtime'=>time(),
                        'cid'=>$user['uid']
                    ));
                    if(!$res){
                        throw new Exception("加资金记录失败12");
                    }

                    //加赠送金额
                    $walletS = $this->db->get("member_wallet",'*',['uid'=>$activeUser['uid'],'coin_id'=>$this->coinId]);
                    if(empty($walletS)){
                        $res = $this->db->insert("member_wallet",array(
                            'uid'=>$activeUser['uid'],
                            'coin_id'=>$this->coinId,
                            'total'=>0,
                            'frozen'=>0
                        ));
                        if(!$res){
                            throw new Exception("加资金账户失败");
                        }
                        $walletS['id'] = $this->db->id();
                    }
                    $res = $this->db->update("member_wallet",['total[+]'=>$num2],['id'=>$walletS['id']]);
                    if(!$res){
                        throw new Exception("减金额失败");
                    }
                    //资金记录
                    $res = $this->db->insert("member_wallet_coin_log",array(
                        'uid'=>$activeUser['uid'],
                        'coin_id'=>$this->coinId,
                        'coin_name'=>$this->coinName,
                        'type'=>103,
                        'note'=>'激活账号赠送',
                        'value'=> $num2,
                        'balance'=>$walletS['total'] ?: 0,
                        'wallet_id'=>$walletS['id'],
                        'status'=>1,
                        'inputtime'=>time(),
                        'cid'=>$user['uid']
                    ));
                    if(!$res){
                        throw new Exception("加资金记录失败1");
                    }
                }
                //添加被激活记录
                $res = $this->db->insert("member_wallet_coin_log",array(
                    'uid'=>$activeUser['uid'],
                    'coin_id'=>$this->coinId,
                    'coin_name'=>$this->coinName,
                    'type'=>108,
                    'note'=>'已被'.substr($user['username'],0,2).'****'.substr($user['username'],-2).'激活',
                    'value'=> 0,
                    'balance'=>0,
                    'wallet_id'=>$walletS['id'],
                    'status'=>1,
                    'inputtime'=>time(),
                    'cid'=>$user['uid']
                ));
                if(!$res){
                    throw new Exception("加资金记录失败3");
                }

                //建立邀请关系
                $res = $this->db->insert("member_invite",array(
                    'uid'=>$user['uid'],
                    'username'=>$user['username'],
                    'rid'=>$activeUser['uid'],
                    'rname'=>$activeUser['username'],
                    'regtime'=>time()
                ));
                if(!$res){
                    throw new Exception("建立邀请关系失败");
                }
                return true;
            }catch (Exception $e){
                return false;
            }
        });
        if($res === true){
            $this->echo_json(1,'激活成功');
        }
        $this->echo_json(0,'激活失败');
    }

    //创建红包
    public function createRedBagAction(){
        $type = intval($this->requestParams['type']);//红包类型 1 分享红包  2 主题红包
        $amount = floatval($this->requestParams['amount']);  //红包金额
        $num = intval($this->requestParams['num']); //红包个数
        $code = trim($this->requestParams['code']); //红包个数
        $deal_pwd = trim($this->requestParams['deal_pwd']); //红包个数
        $topic_id = intval($this->requestParams['topic_id']); //红包个数
        if(!$this->uid || !$deal_pwd){
            $this->echo_json(0,'参数错误');
        }
        //没有实名不能发红包
        $auth = $this->db->get("member_auth",'*',['uid'=>$this->uid,'status'=>3]);
        if(empty($auth)){
            $this->echo_json(0, '请实名后再发红包');
        }
        if($type != 1 && $type != 2){
            $this->echo_json(0,'请选择正确的红包类型');
        }
        if($amount <= 0){
            $this->echo_json(0,'请输入有效的金额');
        }
        if($num <= 0){
            $this->echo_json(0,'请输入有效的数量');
        }
//        if(empty($code)){
//            $this->echo_json(0,'请输入短信验证码');
//        }
        $member = $this->db->get("member",'*',['uid'=>$this->uid]);
        if(empty($member)){
            $this->echo_json(0,'用户不存在');
        }
//        $smsCode = $this->redis->get("sms_".$member['phone']);
//        if($smsCode != $code){
//            $this->echo_json(0,'短信验证码错误');
//        }
        $userinfo = $this->db->get('member_wallet_userinfo','*',['uid'=>$this->uid]);
        if(!$this->password_verify($deal_pwd,$userinfo['deal_psw'], $member['salt'])){
            $this->echo_json(0, '支付密码错误');
        }
        $setting = $this->db->get("member_red_bag_setting",'*');
        if($setting['min'] > 0 && $amount < $setting['min']){
            $this->echo_json(0,'红包金额不得低于'.round($setting['min'],4));
        }
        $wallet = $this->db->get("member_wallet",'*',['uid'=>$this->uid,'coin_id'=>$this->coinId]);
        if(empty($wallet) || $wallet['total'] < ($amount+$setting['fee'])){
            $this->echo_json(0,"余额不足");
        }
        $res = $this->db->action(function() use($wallet,$setting,$amount,$num,$type,$topic_id){
            try{
                $total = $amount + $setting['fee'];
                //红包记录
                $res = $this->db->insert("member_red_bag",array(
                    'uid'=>$this->uid,
                    'amount'=>$amount,
                    'num'=>$num,
                    'left_num'=>$num,
                    'type'=>$type,
                    'topic_id'=>$topic_id,
                    'inputtime'=>time(),
                    'expiration'=>time()+48*3600,
                ));
                if(!$res){
                    throw new Exception('添加红包记录失败');
                }
                $logId = $this->db->id();

                //分配红包  存入表
                $min = 0;//每个人最少能收到0元
                for ($i=1;$i<=$num;$i++)
                {
                    if($i < $num){
                        $safe_total = ($amount-($num-$i)*$min)/($num-$i);//随机安全上限
                        $money = mt_rand($min*100,$safe_total*100)/100;
                        $amount = $amount-$money;
                    }else{
                        $money = $amount;
                    }
                    $res = $this->db->insert("member_red_bag_info",array(
                        'amount'=>$money,
                        'inputtime'=>time(),
                        'bag_id'=>$logId,
                    ));
                    if(!$res){
                        throw new Exception('添加红包记录详情失败');
                    }
                }
                //减金额
                $res = $this->db->update("member_wallet",['total[-]'=>$total],['id'=>$wallet['id']]);
                if(!$res){
                    throw new Exception("减金额失败");
                }
                $note = '';
                if($type == 1){
                    $note = '邀请好友包';
                }
                if($type == 2){
                    $note = '主图红包';
                }
                //资金记录
                $res = $this->db->insert("member_wallet_coin_log",array(
                    'uid'=>$this->uid,
                    'coin_id'=>$this->coinId,
                    'coin_name'=>$this->coinName,
                    'type'=>104,
                    'note'=> $note,
                    'value'=> -1 * ($total - $setting['fee']),
                    'balance'=>$wallet['total'],
                    'wallet_id'=>$wallet['id'],
                    'status'=>1,
                    'inputtime'=>time(),
                    'cid'=>$logId
                ));
                if(!$res){
                    throw new Exception("加资金记录失败");
                }
                if($setting['fee'] > 0){
                    //资金记录
                    $res = $this->db->insert("member_wallet_coin_log",array(
                        'uid'=>$this->uid,
                        'coin_id'=>$this->coinId,
                        'coin_name'=>$this->coinName,
                        'type'=>105,
                        'note'=> '红包手续费',
                        'value'=> -1 * $setting['fee'],
                        'balance'=>  ($wallet['total'] - $amount),
                        'wallet_id'=> $wallet['id'],
                        'status'=>1,
                        'inputtime'=>time(),
                        'cid'=>$logId
                    ));
                    if(!$res){
                        throw new Exception("加资金记录失败1");
                    }
                }
                $GLOBALS['logId'] = $logId;
                return true;
            }catch (Exception $e){
                return false;
            }
        });
        $this->redis->del("sms_".$member['phone']);
        if($res === true){
            $this->echo_json(1,'ok',$GLOBALS['logId']);
        }
        $this->echo_json(0,'失败');
    }

    //红包手续费
    public function redBagFeeAction(){
        $setting = $this->db->get("member_red_bag_setting",'*');
        $this->echo_json(1,'ok',$setting['fee'] ?: 0);
    }
    //抢红包
    public function grabRedEnvelopeAction(){
//        $username = trim($this->requestParams['username']);
//        $password = trim($this->requestParams['password']);
        $userRed = trim($this->requestParams['user_red']);
        $redBagId = intval($this->requestParams['id']);
//        if(!$username || !$password || !$userRed || !$redBagId){
//            $this->echo_json(0,'请填写完整信息');
//        }
        if(!$this->uid || !$userRed || !$redBagId){
            $this->echo_json(0,'请填写完整信息1');
        }
        $user = $this->db->get("member",'*',['uid'=>$this->uid]);
        if(empty($user)){
            $this->echo_json(0,'用户不存在');
        }
        //验证密码
//        if (!$this->password_verify($password, $user['password'], $user['salt'])) {
//            $this->echo_json(0, '密码错误');
//        }
        //发红包用户是否存在
        $redUser = $this->db->get("member",'*',['username'=>$userRed]);
        if(empty($redUser)){
            $this->echo_json(0, '发红包用户不存在');
        }
        //红包是否存在
        $redBag = $this->db->get("member_red_bag",'*',['uid'=>$redUser['uid'],'type'=>2,'left_num[>]'=>0,'id'=>$redBagId,'expiration[>]'=>time()]);
        if(!$redBag){
            $this->echo_json(0,'红包已经领完了1');
        }
        //当前用户是否已经领过当前红包
        $current = $this->db->get("member_red_bag_info",'*',['bag_id'=>$redBag['id'],'uid'=>$this->uid]);
        if($current){
            $this->echo_json(1,'ok');
        }
        $redInfo = $this->db->select("member_red_bag_info",'*',['bag_id'=>$redBag['id'],'status'=>1]);
        if(!$redInfo){
            $this->echo_json(0,'红包已经领完了');
        }
        //随机抽取一个红包
        $red = $redInfo[array_rand($redInfo,1)];
        $res = $this->db->action(function() use($redInfo,$redBag,$user,$red){
            try{
                //减红包个数
                $res = $this->db->update("member_red_bag",['left_num[-]'=>1],['id'=>$redBag['id']]);
                if(!$res){
                    throw new Exception('减红包个数失败');
                }
                //加红包金额
                if($red['amount'] > 0){
                    $wallet = $this->db->get("member_wallet",'*',['uid'=>$user['uid'],'coin_id'=>$this->coinId]);
                    if(!$wallet){
                        throw new Exception("资金账户不存在");
                    }
                    $res = $this->db->update("member_wallet",['total[+]'=>$red['amount']],['id'=>$wallet['id']]);
                    if(!$res){
                        throw new Exception("加金额失败");
                    }
                    //资金记录
                    $res = $this->db->insert("member_wallet_coin_log",array(
                        'uid'=>$user['uid'],
                        'coin_id'=>$this->coinId,
                        'coin_name'=>$this->coinName,
                        'type'=>106,
                        'note'=> '领取主题红包',
                        'value'=> $red['amount'],
                        'balance'=>$wallet['total'],
                        'wallet_id'=>$wallet['id'],
                        'status'=>1,
                        'inputtime'=>time(),
                        'cid'=>$red['id']
                    ));
                    if(!$res){
                        throw new Exception("加资金记录失败");
                    }
                }
                //修改红包状态
                $res = $this->db->update("member_red_bag_info",['status[-]'=>1,'updatetime'=>time(),'uid'=>$user['uid']],['id'=>$red['id']]);
                if(!$res){
                    throw new Exception("修改红包状态失败");
                }
                return true;
            }catch (Exception $e){
                return false;
            }
        });
        if($res === true){
            $this->echo_json(1,'ok');
        }
        $this->echo_json(0,'领取红包失败');
    }

    //领取完红包页面
    public function redDoneAction(){
        $id = intval($this->requestParams['id']);
        $page = intval($this->requestParams['page']) ?: 1;
        if(!$id || !$this->uid){
            $this->echo_json(0,'参数错误');
        }
        $bag = $this->db->get("member_red_bag",'*',['id'=>$id]);
        if(!$bag){
            $this->echo_json(0,'红包信息错误');
        }
        $list = $this->db->select("member_red_bag_info",['[>]member'=>['uid'=>'uid']],[
            'member_red_bag_info.amount',
            'member_red_bag_info.updatetime',
            'member.username',
            'member.uid',
        ],
            [
                'member_red_bag_info.status'=>0,
                'member_red_bag_info.bag_id'=>$bag['id'],
                'LIMIT'=>[($page -1)*10,10],
                'ORDER'=>['member_red_bag_info.updatetime'=>'DESC']
            ]);
        $num = 0;
        foreach ($list as $k=>$v){
            if($v['uid'] == $this->uid){
                $num = $v['amount'];
            }
            $list[$k]['updatetime'] = date("Y/m/d H:i:s",$v['updatetime']);
            $list[$k]['username'] = substr($v['username'],0,1).'***'.substr($v['username'],-1);
        }
        $ret = [
            'num'=>$num,
            'coinName'=>$this->coinName,
            'data'=>$bag,
            'list'=>$list
        ];
        $this->echo_json(1,'ok',$ret);
    }

    //根据用户名查找用户
    public function memberInfoAction(){
        $username = trim($this->requestParams['username']);
        if(!$username){
            $this->echo_json(1,'ok');
        }
        $user = $this->db->get("member",'*',['username'=>$username]);
        if(!$user){
            $this->echo_json(1,'ok');
        }
        $data['username'] = $user['username'];
        $data['avatar'] = $user['avatar'];
        $data['phone'] = $user['phone'];
        $data['area'] = $user['area'];
        $data['googlecode'] = $user['googlecode'] ?: 0;
        $data['googleoff'] = $user['googleoff'] ?: 0;
        $data['coinName'] = $this->coinName;
        $this->echo_json(1,'ok',$data);
    }

    //注册
    public function registerAction()
    {
        $data = $this->requestParams['data'];
        $data['username'] = trim($data['username']);
        $data['phone'] = trim($data['phone']);
        $data['password'] = trim($data['password']);
        $data['password2'] = trim($data['password2']);
        $data['deal_psw'] = trim($data['deal_psw']);
        $data['deal_psw2'] = trim($data['deal_psw2']);
        $data['invite'] = isset($data['invite']) ? trim($data['invite']) : '';
        $data['email'] = trim($data['email']);
        $data['regtype'] = trim($data['regtype'])?:1;
        $data['red_id'] = intval($data['red_id']);
        $data['red_username'] = trim($data['red_username']);
        if (empty($data['username']) || empty($data['password']) || empty($data['password2'])
            || empty($data['deal_psw'] || empty($data['phone']))
        ) {
            $this->echo_json(0, '账号，密码，确认密码，支付密码必填',$data);
        }
        //判断短信验证码
        if (empty($data['code'])) {
            $this->echo_json(0, '验证码必填');
        }
//        if (1 == $data['regtype']) {
//            $smsuid = $this->redis->get('sms_' . $data['phone']);
////            $smsuid = $this->redis->get('sms_' . $data['username']);
//        } else {
//            $smsuid = $this->redis->get('email_' . $data['username']);
//        }

//        if($code != $data['code']){
//            //验证码错误
//            $this->echo_json(0, '验证码错误');
//        }

//        $smsuid = $this->redis->get('sms_' . $data['phone']);
        $smsuid = $this->redis->get('sms_' . $data['phone']);
        $this->redis->del('sms_' . $data['phone']);//删除验证码
        if (empty($smsuid) && !empty($data['email'])) {
            $smsuid = $this->redis->get('email_' . $data['email']);
            $this->redis->del('email_' . $data['email']);//删除验证码
        }

        if ($smsuid != $data['code']) {
            if ($this->config->application->sms->currencyviticode != $data['code']) {
                $this->echo_json(0, '验证码错误');
            }
        }
        $data['area'] = isset($data['area']) ? $data['area'] : "86";
//        if(!preg_match("/(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{3,16}/",$data['username']) || strlen($data['username']) < 3 || strlen($data['username']) > 16){
//            $this->echo_json(0,'请输入3-16位字母+数字组成的用户名');
//        }
        if(strlen($data['username']) < 3 || strlen($data['username']) > 16){
            $this->echo_json(0,'请输入3-16位字母数字的用户名');
        }
        if(strlen($data['password']) < 6 || strlen($data['password']) > 16 || !preg_match("/(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}/",$data['password'])){
            $this->echo_json(0,'请输入6-16位字母+数字组成的登录密码');
        }
        if(strlen($data['deal_psw']) !=  6 || !preg_match("/\d{6}/",$data['deal_psw'])){
            $this->echo_json(0,'请输入6位纯数字组成的支付密码');
        }
        if ($data['password'] != $data['password2']) {
            $this->echo_json(0, '两次输入的密码不一致');
        }

        if ($data['deal_psw'] != $data['deal_psw2']) {
            $this->echo_json(0, '两次输入的支付密码不一致');
        }

        //判断账号是否重复
        $member = $this->db->get('member', '*', ['username' => $data['username']]);
        if (!empty($member)) {
            $this->echo_json(0, '账号已存在');
        }
//        if(empty($data['invite'])) {
//            $this->echo_json(0, '请填写邀请码');
//        }
        $redBag = $redInfo = [];
        if($data['red_username'] && $data['red_id']){
            $red_user = $this->db->get("member",'*',['username'=>$data['red_username']]);
            if($red_user){
                $redBag = $this->db->get("member_red_bag",'*',['type'=>1,'left_num[>]'=>0,'id'=>$data['red_id'],'uid'=>$red_user['uid'],'expiration[>]'=>time()]);
                if($redBag){
                    $redInfo = $this->db->select("member_red_bag_info",'*',['bag_id'=>$redBag['id'],'status'=>1]);
                }
            }
        }
        $ret = $this->db->action(function ($db) use ($data,$redBag,$redInfo) {
            try {
                // 正常注册时，会员初始化信息
                $salt = '';
                $chars = '0123456789abcdefghijklmnopqrstuvwxyz';
                $max = strlen($chars) - 1;
                mt_srand((double)microtime() * 1000000);
                for ($i = 0; $i < 22; $i++) {
                    $salt .= $chars[mt_rand(0, $max)];
                }
                $spassword = @$this->password_hash($data['password'], $salt);
                $regip = "0";
                $token = "";
                $groupid = 3;
                $randcode = rand(100000, 999999);
                $memberdata = [
                    'salt' => $salt,
                    'name' => '',
                    'email' => $data['email'],
                    'regip' => $regip,
                    'avatar' => '',
                    'regtime' => time(),
                    'groupid' => $groupid,
                    'levelid' => 0,
                    'overdue' => 0,
                    'username' => $data['username'],
                    'password' => $spassword,
                    'randcode' => $randcode,
                    'ismobile' => 0,
                    'token' => $token,
                    'phone' => $data['phone'],
                    'status' => 2,
                    'area' => $data['area'],
                    'addr' => md5($data['username'].$data['phone']).$salt,
                ];
                if (strpos($data['username'],'@') > 0) {
                    $memberdata['email'] = $data['username'];
                } else {
//                    $memberdata['phone'] = $data['username'];
                }
                $me = $db->insert('member', $memberdata);
                if (!$me) {
                    throw new Exception("注册失败1");
                }

                $uid = $db->id();
                if ($data['username'] == 'null') {
                    // 防止重名
                    $data['username'] = $uid;
                    $db->update('member', [
                        'username' => $data['username'],
                    ], ['uid' => $uid]);
                }

                $member_data['uid'] = $uid;

                // 邀请注册
                $char = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';

                $yaoqingma = '';
                for ($i = 0; $i < 10; $i++) {
                    $yaoqingma .= $char[mt_rand(0, 53)];
                }
                $yaoqingmaright = '';
                for ($i = 0; $i < 10; $i++) {
                    $yaoqingmaright .= $char[mt_rand(0, 53)];
                }
                $member_data['yaoqingma'] = $yaoqingma;

                //邀请码 dr_member_data
                if (!empty($member_data['yaoqingma'])) {

                    $re = $db->insert('member_data', [
                        'uid' => $uid,
                        'yaoqingma' => $uid.$member_data['yaoqingma'],
                        'area' => $data['area'],
                        'yaoqingmaright' => $uid.$yaoqingmaright,
                    ]);
                    if (!$re) {
                        throw new Exception("邀请码添加失败");
                    }
                }
//                if (!empty($data['invite'])) {
//                    $memberinfo = $db->get('member_data', '*', ['OR'=>['yaoqingma' => $data['invite'],'yaoqingmaright'=>$data['invite']]]);
//                    $parents = $memberinfo['parents'];
//                    if (empty($memberinfo)) {
//                        $memberinfo = $db->get('member', '*', ['uid' => $data['invite']]);
//                        if (empty($memberinfo)) {
//                            $this->echo_json(0, '邀请人不存在');
//                            //throw new Exception("邀请人不存在");
//                        }
//                        $member_data =  $db->get('member_data', '*',array('uid'=>$memberinfo['uid']));
//                        $parents = $member_data['parents'];
//                    }
//
//                    $member_uid = $db->get('member', '*', ['uid' => $memberinfo['uid']]);
//
//                    $idata = [
//                        'uid' => $memberinfo['uid'],
//                        'rid' => $uid,
//                        'rname' => $data['username'],
//                        'regtime' => time(),
//                        'username' => $member_uid['username']
//                    ];
//                    $re = $this->db->insert('member_invite', $idata);
//
//                    if (!$re) {
//                        throw new Exception("失败");
//                    }
//                    $db->update('member_data',array(
//                        'parents'=>$parents.$memberinfo['uid'].','
//                    ),array('uid'=>$uid));
//                }else{
//                    $db->update('member_data',array('parents'=>','),array('uid'=>$uid));
//                }
//                if (!empty($data['invite'])) {
//                    $pmemberinfo = $db->get('member_data','*',['yaoqingmaright'=>$data['invite']]);
//                    if (!empty($pmemberinfo)) {
//                        $pwpuid = $this->dicz($pmemberinfo['uid'], 2, $uid);
//                        $db->update('member_invite',['pwpuid'=>$pwpuid,'type'=>2],['rid'=>$uid]);
//                    } else {
//                        $pmemberinfo = $db->get('member_data','*',['yaoqingma'=>$data['invite']]);
//                        $pwpuid = $this->dicz($pmemberinfo['uid'], 1, $uid);
//                        $db->update('member_invite',['pwpuid'=>$pwpuid],['rid'=>$uid]);
//                    }
//                }
                //钱包处理
                $list = $db->select('member_wallet_coin','*');
                foreach($list as $k => $v){
                    //添加钱包
                    $re = $db->insert('member_wallet',[
                        'uid' => $uid,
                        'coin_id' => $v['id'],
                        'total' => 0,
                        'frozen' => 0,
                        'inputtime' => date('Y-m-d H:i:s'),
                        'updatetime' => date('Y-m-d H:i:s')
                    ]);
                    if($re <= 0) {
                        throw new Exception('钱包处理失败');
                    }
                }
                //生成私钥
                $char = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';
                $pkey = '';
                for ($i = 0; $i < 64; $i++) {
                    $pkey .= $char[mt_rand(0, 53)];
                }
                //生成助记词
                //            $list = $this->db->select(' *,RAND() as r ')->limit(10)->order_by('r')
                //                ->get('member_wallet_words')->result_array();
                $lists = $db->query("select *,RAND() as r  from dr_member_wallet_words  order by r LIMIT 10");
                $list = $lists->fetchAll();
                $word = '';
                foreach ($list as $k => $v) {
                    $word .= $v['word'] . ',';
                }
                //支付密码  member_wallet_userinfo
                if (empty($data['deal_psw'])) {
                    $re = $db->insert('member_wallet_userinfo', [
                            'uid' => $uid,
                            'pkey' => $pkey,
                            'word' => substr($word, 0, -1),
                            'deal_psw' => ''
                        ]
                    );
                    if (!$re) {
                        throw new Exception('支付密码设置失败');
                    }
                } else {

                    if ($data['password'] == $data['deal_psw']) {
                        $this->echo_json('0', '支付密码不能和登录密码相同');
                    }

                    $re = $db->insert('member_wallet_userinfo', array(
                            'uid' => $uid,
                            'pkey' => $pkey,
                            'word' => substr($word, 0, -1),
                            'deal_psw' => $this->password_hash($data['deal_psw'], $salt),
                        )
                    );
                    if (!$re) {
                        throw new Exception('支付密码设置失败1');
                    }
                }
                //领取红包
                if($data['red_id'] && $data['red_username'] && $redBag && $redInfo){
                    //红包是否存在
                    //随机抽取一个红包
                    $red = $redInfo[array_rand($redInfo,1)];
                    //减红包个数
                    $res = $this->db->update("member_red_bag",['left_num[-]'=>1],['id'=>$redBag['id']]);
                    if(!$res){
                        throw new Exception('减红包个数失败');
                    }
                    $GLOBALS['redNums'] = $red['amount'];
                    //加红包金额
                    if($red['amount'] > 0){
                        $wallet = $this->db->get("member_wallet",'*',['uid'=>$uid,'coin_id'=>$this->coinId]);
                        if(!$wallet){
                            throw new Exception("资金账户不存在");
                        }
                        $res = $this->db->update("member_wallet",['total[+]'=>$red['amount']],['id'=>$wallet['id']]);
                        if(!$res){
                            throw new Exception("加金额失败");
                        }
                        //资金记录
                        $res = $this->db->insert("member_wallet_coin_log",array(
                            'uid'=>$uid,
                            'coin_id'=>$this->coinId,
                            'coin_name'=>$this->coinName,
                            'type'=>106,
                            'note'=> '领取邀请红包',
                            'value'=> $red['amount'],
                            'balance'=>$wallet['total'],
                            'wallet_id'=>$wallet['id'],
                            'status'=>1,
                            'inputtime'=>time(),
                            'cid'=>$red['id']
                        ));
                        if(!$res){
                            throw new Exception("加资金记录失败");
                        }
                    }
                    //修改红包状态
                    $res = $this->db->update("member_red_bag_info",['status[-]'=>1,'updatetime'=>time(),'uid'=>$uid],['id'=>$red['id']]);
                    if(!$res){
                        throw new Exception("修改红包状态失败");
                    }
                }
                return true;
            } catch (Exception $exception) {
                return false;
            }
        });
        if (1 == $data['regtype']) {
//            $this->redis->del('sms_' . $data['username']);//删除验证码
            $this->redis->del('sms_' . $data['phone']);//删除验证码
        } else {
            $this->redis->del('email_' . $data['username']);//删除验证码
        }
        if ($ret === true) {
            $this->echo_json(1, '注册成功',['num'=>$GLOBALS['redNums'],'coinName'=>$this->coinName]);
        }else{
            $this->echo_json('0', '系统繁忙');
        }
    }

    public function registerTestAction()
    {
        $data = $this->requestParams['data'];
        $data['username'] = trim($data['username']);
        $data['phone'] = trim($data['phone']);
        $data['password'] = trim($data['password']);
        $data['password2'] = trim($data['password2']);
        $data['deal_psw'] = trim($data['deal_psw']);
        $data['deal_psw2'] = trim($data['deal_psw2']);
        $data['invite'] = isset($data['invite']) ? trim($data['invite']) : '';
        $data['email'] = trim($data['email']);
        $data['regtype'] = trim($data['regtype'])?:1;
        $data['red_id'] = intval($data['red_id']);
        $data['red_username'] = trim($data['red_username']);
        if (empty($data['username']) || empty($data['password']) || empty($data['password2'])
            || empty($data['deal_psw'] || empty($data['phone']))
        ) {
            $this->echo_json(0, '账号，密码，确认密码，支付密码必填',$data);
        }
        //判断短信验证码
        if (empty($data['code'])) {
            $this->echo_json(0, '验证码必填');
        }
//        if (1 == $data['regtype']) {
//            $smsuid = $this->redis->get('sms_' . $data['phone']);
////            $smsuid = $this->redis->get('sms_' . $data['username']);
//        } else {
//            $smsuid = $this->redis->get('email_' . $data['username']);
//        }

        $code = $this->redis->get('sms_' . $data['phone']);
        $this->redis->del('sms_' . $data['phone']);//删除验证码
        if (empty($code) && !empty($data['email'])) {
            $code = $this->redis->get('email_' . $data['email']);
            $this->redis->del('email_' . $data['email']);//删除验证码
        }
//        if($code != $data['code']){
//            //验证码错误
//            $this->echo_json(0, '验证码错误');
//        }
//        $smsuid = $this->redis->get('sms_' . $data['phone']);
        if ($code != $data['code']) {
            if ($this->config->application->sms->currencyviticode != $data['code']) {
                $this->echo_json(0, '验证码错误');
            }
        }
        $data['area'] = isset($data['area']) ? $data['area'] : "86";
//        if(!preg_match("/(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{3,16}/",$data['username']) || strlen($data['username']) < 3 || strlen($data['username']) > 16){
//            $this->echo_json(0,'请输入3-16位字母+数字组成的用户名');
//        }
        if(strlen($data['username']) < 3 || strlen($data['username']) > 16){
            $this->echo_json(0,'请输入3-16位字母数字的用户名');
        }
        if(strlen($data['password']) < 6 || strlen($data['password']) > 16 || !preg_match("/(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}/",$data['password'])){
            $this->echo_json(0,'请输入6-16位字母+数字组成的登录密码');
        }
        if(strlen($data['deal_psw']) !=  6 || !preg_match("/\d{6}/",$data['deal_psw'])){
            $this->echo_json(0,'请输入6位纯数字组成的支付密码');
        }
        if ($data['password'] != $data['password2']) {
            $this->echo_json(0, '两次输入的密码不一致');
        }

        if ($data['deal_psw'] != $data['deal_psw2']) {
            $this->echo_json(0, '两次输入的支付密码不一致');
        }

        //判断账号是否重复
        $member = $this->db->get('member', '*', ['username' => $data['username']]);
        if (!empty($member)) {
            $this->echo_json(0, '账号已存在');
        }
//        if(empty($data['invite'])) {
//            $this->echo_json(0, '请填写邀请码');
//        }
        $redBag = $redInfo = [];
        if($data['red_username'] && $data['red_id']){
            $red_user = $this->db->get("member",'*',['username'=>$data['red_username']]);
            if($red_user){
                $redBag = $this->db->get("member_red_bag",'*',['type'=>1,'left_num[>]'=>0,'id'=>$data['red_id'],'uid'=>$red_user['uid'],'expiration[>]'=>time()]);
                if($redBag){
                    $redInfo = $this->db->select("member_red_bag_info",'*',['bag_id'=>$redBag['id'],'status'=>1]);
                }
            }
        }
        $ret = $this->db->action(function ($db) use ($data,$redBag,$redInfo) {
            try {
                // 正常注册时，会员初始化信息
                $salt = '';
                $chars = '0123456789abcdefghijklmnopqrstuvwxyz';
                $max = strlen($chars) - 1;
                mt_srand((double)microtime() * 1000000);
                for ($i = 0; $i < 22; $i++) {
                    $salt .= $chars[mt_rand(0, $max)];
                }
                $spassword = @$this->password_hash($data['password'], $salt);
                $regip = "0";
                $token = "";
                $groupid = 3;
                $randcode = rand(100000, 999999);
                $memberdata = [
                    'salt' => $salt,
                    'name' => '',
                    'email' => $data['email'],
                    'regip' => $regip,
                    'avatar' => '',
                    'regtime' => time(),
                    'groupid' => $groupid,
                    'levelid' => 0,
                    'overdue' => 0,
                    'username' => $data['username'],
                    'password' => $spassword,
                    'randcode' => $randcode,
                    'ismobile' => 0,
                    'token' => $token,
                    'phone' => $data['phone'],
                    'status' => 2,
                    'area' => $data['area'],
                    'addr' => md5($data['username'].$data['phone']).$salt,
                ];
                if (strpos($data['username'],'@') > 0) {
                    $memberdata['email'] = $data['username'];
                } else {
//                    $memberdata['phone'] = $data['username'];
                }
                $me = $db->insert('member', $memberdata);
                if (!$me) {
                    throw new Exception("注册失败1");
                }

                $uid = $db->id();
                if ($data['username'] == 'null') {
                    // 防止重名
                    $data['username'] = $uid;
                    $db->update('member', [
                        'username' => $data['username'],
                    ], ['uid' => $uid]);
                }

                $member_data['uid'] = $uid;

                // 邀请注册
                $char = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';

                $yaoqingma = '';
                for ($i = 0; $i < 10; $i++) {
                    $yaoqingma .= $char[mt_rand(0, 53)];
                }
                $yaoqingmaright = '';
                for ($i = 0; $i < 10; $i++) {
                    $yaoqingmaright .= $char[mt_rand(0, 53)];
                }
                $member_data['yaoqingma'] = $yaoqingma;

                //邀请码 dr_member_data
                if (!empty($member_data['yaoqingma'])) {

                    $re = $db->insert('member_data', [
                        'uid' => $uid,
                        'yaoqingma' => $uid.$member_data['yaoqingma'],
                        'area' => $data['area'],
                        'yaoqingmaright' => $uid.$yaoqingmaright,
                    ]);
                    if (!$re) {
                        throw new Exception("邀请码添加失败");
                    }
                }
//                if (!empty($data['invite'])) {
//                    $memberinfo = $db->get('member_data', '*', ['OR'=>['yaoqingma' => $data['invite'],'yaoqingmaright'=>$data['invite']]]);
//                    $parents = $memberinfo['parents'];
//                    if (empty($memberinfo)) {
//                        $memberinfo = $db->get('member', '*', ['uid' => $data['invite']]);
//                        if (empty($memberinfo)) {
//                            $this->echo_json(0, '邀请人不存在');
//                            //throw new Exception("邀请人不存在");
//                        }
//                        $member_data =  $db->get('member_data', '*',array('uid'=>$memberinfo['uid']));
//                        $parents = $member_data['parents'];
//                    }
//
//                    $member_uid = $db->get('member', '*', ['uid' => $memberinfo['uid']]);
//
//                    $idata = [
//                        'uid' => $memberinfo['uid'],
//                        'rid' => $uid,
//                        'rname' => $data['username'],
//                        'regtime' => time(),
//                        'username' => $member_uid['username']
//                    ];
//                    $re = $this->db->insert('member_invite', $idata);
//
//                    if (!$re) {
//                        throw new Exception("失败");
//                    }
//                    $db->update('member_data',array(
//                        'parents'=>$parents.$memberinfo['uid'].','
//                    ),array('uid'=>$uid));
//                }else{
//                    $db->update('member_data',array('parents'=>','),array('uid'=>$uid));
//                }
//                if (!empty($data['invite'])) {
//                    $pmemberinfo = $db->get('member_data','*',['yaoqingmaright'=>$data['invite']]);
//                    if (!empty($pmemberinfo)) {
//                        $pwpuid = $this->dicz($pmemberinfo['uid'], 2, $uid);
//                        $db->update('member_invite',['pwpuid'=>$pwpuid,'type'=>2],['rid'=>$uid]);
//                    } else {
//                        $pmemberinfo = $db->get('member_data','*',['yaoqingma'=>$data['invite']]);
//                        $pwpuid = $this->dicz($pmemberinfo['uid'], 1, $uid);
//                        $db->update('member_invite',['pwpuid'=>$pwpuid],['rid'=>$uid]);
//                    }
//                }
                //钱包处理
                $list = $db->select('member_wallet_coin','*');
                foreach($list as $k => $v){
                    //添加钱包
                    $re = $db->insert('member_wallet',[
                        'uid' => $uid,
                        'coin_id' => $v['id'],
                        'total' => 0,
                        'frozen' => 0,
                        'inputtime' => date('Y-m-d H:i:s'),
                        'updatetime' => date('Y-m-d H:i:s')
                    ]);
                    if($re <= 0) {
                        throw new Exception('钱包处理失败');
                    }
                }
                //生成私钥
                $char = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';
                $pkey = '';
                for ($i = 0; $i < 64; $i++) {
                    $pkey .= $char[mt_rand(0, 53)];
                }
                //生成助记词
                //            $list = $this->db->select(' *,RAND() as r ')->limit(10)->order_by('r')
                //                ->get('member_wallet_words')->result_array();
                $lists = $db->query("select *,RAND() as r  from dr_member_wallet_words  order by r LIMIT 10");
                $list = $lists->fetchAll();
                $word = '';
                foreach ($list as $k => $v) {
                    $word .= $v['word'] . ',';
                }
                //支付密码  member_wallet_userinfo
                if (empty($data['deal_psw'])) {
                    $re = $db->insert('member_wallet_userinfo', [
                            'uid' => $uid,
                            'pkey' => $pkey,
                            'word' => substr($word, 0, -1),
                            'deal_psw' => ''
                        ]
                    );
                    if (!$re) {
                        throw new Exception('支付密码设置失败');
                    }
                } else {

                    if ($data['password'] == $data['deal_psw']) {
                        $this->echo_json('0', '支付密码不能和登录密码相同');
                    }

                    $re = $db->insert('member_wallet_userinfo', array(
                            'uid' => $uid,
                            'pkey' => $pkey,
                            'word' => substr($word, 0, -1),
                            'deal_psw' => $this->password_hash($data['deal_psw'], $salt),
                        )
                    );
                    if (!$re) {
                        throw new Exception('支付密码设置失败1');
                    }
                }
                //领取红包
                if($data['red_id'] && $data['red_username'] && $redBag && $redInfo){
                    //红包是否存在
                    //随机抽取一个红包
                    $red = $redInfo[array_rand($redInfo,1)];
                    //减红包个数
                    $res = $this->db->update("member_red_bag",['left_num[-]'=>1],['id'=>$redBag['id']]);
                    if(!$res){
                        throw new Exception('减红包个数失败');
                    }
                    $GLOBALS['redNums'] = $red['amount'];
                    //加红包金额
                    if($red['amount'] > 0){
                        $wallet = $this->db->get("member_wallet",'*',['uid'=>$uid,'coin_id'=>$this->coinId]);
                        if(!$wallet){
                            throw new Exception("资金账户不存在");
                        }
                        $res = $this->db->update("member_wallet",['total[+]'=>$red['amount']],['id'=>$wallet['id']]);
                        if(!$res){
                            throw new Exception("加金额失败");
                        }
                        //资金记录
                        $res = $this->db->insert("member_wallet_coin_log",array(
                            'uid'=>$uid,
                            'coin_id'=>$this->coinId,
                            'coin_name'=>$this->coinName,
                            'type'=>106,
                            'note'=> '领取邀请红包',
                            'value'=> $red['amount'],
                            'balance'=>$wallet['total'],
                            'wallet_id'=>$wallet['id'],
                            'status'=>1,
                            'inputtime'=>time(),
                            'cid'=>$red['id']
                        ));
                        if(!$res){
                            throw new Exception("加资金记录失败");
                        }
                    }
                    //修改红包状态
                    $res = $this->db->update("member_red_bag_info",['status[-]'=>1,'updatetime'=>time(),'uid'=>$uid],['id'=>$red['id']]);
                    if(!$res){
                        throw new Exception("修改红包状态失败");
                    }
                }
                return true;
            } catch (Exception $exception) {
                return false;
            }
        });
        if (1 == $data['regtype']) {
//            $this->redis->del('sms_' . $data['username']);//删除验证码
            $this->redis->del('sms_' . $data['phone']);//删除验证码
        } else {
            $this->redis->del('email_' . $data['username']);//删除验证码
        }
        if ($ret === true) {
            $this->echo_json(1, '注册成功',['num'=>$GLOBALS['redNums'],'coinName'=>$this->coinName]);
        }else{
            $this->echo_json('0', '系统繁忙');
        }
    }
    //生成随机字符串
    public function randChars($num){
        $salt = '';
        $chars = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';
        $max = strlen($chars) - 1;
        mt_srand((double)microtime() * 1000000);
        for ($i = 0; $i < $num; $i++) {
            $salt .= $chars[mt_rand(0, $max)];
        }
        return $salt;
    }
    //登录
    public function loginAction()
    {
        $username = trim($this->requestParams['username']);
        $password = trim($this->requestParams['password']);
        if (empty($username) || empty($password)) {
            $this->echo_json(0, '参数错误');
        }
        $memberInfo = $this->db->get('member', '*', ['username' => $username]);
        if (!is_array($memberInfo) || empty($memberInfo)) {
            $this->echo_json(0, '账户不存在');
        }

        if (!$this->password_verify($password, $memberInfo['password'], $memberInfo['salt'])) {
            $this->echo_json(0, '密码不正确');
        }
        if($memberInfo['status'] == 0) {
            $this->echo_json(0, '该账户已冻结，请联系客服解冻');
        }
        if($memberInfo['status'] == 2) {
            $this->echo_json(400, '该账户未激活');
        }
        //登录其他设备清除前token
        if(!empty($memberInfo['token'])){
            $this->redis->del($memberInfo['token']);
        }

        unset($memberInfo['password']);
        unset($memberInfo['token']);
        $char = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';
        $token = '';
        for($i = 0; $i < 32; $i++) {
            $token .= $char[mt_rand(0, 53)];
        }
        $memberInfo['token'] = $token;
        $this->redis->set($memberInfo['token'], $memberInfo['uid'],1800);
        $this->db->update('member', ['token' => $memberInfo['token']], ['uid' => $memberInfo['uid']]);
        $yaoqm = $this->db->get('member_data',['yaoqingma','yaoqingmaright','area'],['uid'=>$memberInfo['uid']]);
        $this->initAccount($memberInfo['uid']);
        $this->echo_json(1, 'ok',array_merge($memberInfo, $yaoqm));
    }

    //忘记密码
    public function forgetPasswordAction()
    {
        $data = $this->requestParams['data'];
        if ($data['type'] == 1) {
            //电话号码找回
            if (empty($data['phone']) || empty($data['username'])) {
                $this->echo_json(0, '请输入电话号码');
            }
        } elseif($data['type'] == 2) {
            //私钥找回
            if (empty($data['private_key'])) {
                $this->echo_json(0, '请输入私钥');
            }
        }else{
            //邮箱找回
            if (empty($data['email'])) {
                $this->echo_json(0, '请输入邮箱');
            }
        }
        if(strlen($data['password']) < 6 || strlen($data['password']) > 16 || !preg_match("/(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}/",$data['password'])){
            $this->echo_json(0,'请输入6-16位字母+数字组成的登录密码');
        }
        if(strlen($data['password2']) < 6 || strlen($data['password2']) > 16 || !preg_match("/(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}/",$data['password2'])){
            $this->echo_json(0,'请输入6-16位字母+数字组成的登录密码');
        }
//        if (empty($data['password'])) {
//            $this->echo_json(0, '请输入登录密码');
//        }
//        if (empty($data['password2'])) {
//            $this->echo_json(0, '请输入确认密码');
//        }
        if ($data['password'] != $data['password2']) {
            $this->echo_json(0, '两次输入的密码不一致');
        }

        if ($data['type'] == 1) {

            //电话号码找回
            $phone_code = $this->redis->get('sms_' . $data['phone']);
            if (empty($data['code'])) {
                $this->echo_json(0, '请输入验证码');
            }
            if ($phone_code != $data['code']) {
                //验证码不正确
                $this->echo_json(0, '验证码错误');
            }

            $member = $this->db->get('member', '*', ['username' => $data['username'],'phone'=>$data['phone']]);

        } elseif($data['type'] == 2) {
            //私钥找回
            $pkey = $this->db->get('member_wallet_userinfo', '*', ['pkey' => $data['private_key']]);
            if (empty($pkey)) {
                //私钥错误
                $this->echo_json(0, '私钥错误');
            }

            $member = $this->db->get('member', '*', ['uid' => $pkey['uid']]);;
        }else{
            //邮箱找回
            $code = $this->redis->get('email_' . $data['email']);
            if (empty($data['code'])) {
                $this->echo_json(0, '请输入验证码');
            }
            if ($code != $data['code']) {
                //验证码不正确
                $this->echo_json(0, '验证码错误');
            }

            $member = $this->db->get('member', '*', ['email' => $data['email']]);

        }

        if (empty($member)) {
            $this->echo_json(0, '账户未注册');
        }

        $googlecode = $this->requestParams['googlecode'];

        if(!$this->viteCode($googlecode, 4, $member['uid'])) {
            $this->echo_json(0, 'Google验证码错误');
        }

        if ($data['type'] == 1) {
            //删除验证码
            $this->redis->del('sms_' . $data['phone']);
        } else {
            //删除验证码
            $this->redis->del('email_' . $data['email']);
        }

        //修改密码
        //$spassword = $this->db->get('member_wallet_userinfo','*',['uid'=>$uid]);查询支付密码
//        $res = $this->password_verify($data['password'], $spassword['deal_psw']);
//
//        if ($res == true) {
//            $this->echo_json(0, '登录密码不能和支付密码相同');
//        }

        //修改密码

        $spassword = $this->password_hash($data['password'], $member['salt']);
        $this->db->update('member', [
            'password' => $spassword], ['uid' => $member['uid']]
        );
        $this->echo_json(1, 'ok');
    }

    //矿池 用户累计总收益   最近一周收益  往前查7天
    public function minerAction(){
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        //用户矿池总收益
        $minerIns = $this->db->get("member_wallet",'*',['uid'=>$this->uid,'coin_id'=>$this->coinId]);
        //今天最佳持币设置
        $setting = $this->db->get("coin_circulation",'*',['inputtime[>=]'=>strtotime(date('Y-m-d')),'inputtime[<=]'=>time()]);
        //收益记录
        $endTime = time();
        $startTime = strtotime(date("Y-m-d H:i:s",strtotime("-7 days")));
        $list = $this->db->select("member_wallet_coin_log",'*',['uid'=>$this->uid,'type'=>107,'inputtime[>=]'=>$startTime,'inputtime[<=]'=>$endTime]);
        $arr = [];
        foreach ($list as $k=>$v){
            $arr[date("Y-m-d",$v['inputtime'])]['value'] += $v['value'];
            $arr[date("Y-m-d",$v['inputtime'])]['inputtime'] = $v['inputtime'];
        }
        $data = [];
        $i = 0;
        foreach ($arr as $key=>$val){
            $data[$i]['inputtime'] = $val['inputtime'];
            $data[$i]['value'] = number_format($val['value'],6,'.','');
            $i++;
        }
        $ret = array(
            'total'=>number_format($minerIns['miner_total'],6,'.',''),
            'min'=> $this->db->get("setting",'content',['type'=>8]) ?: 0,
            'best'=> number_format($this->db->get("setting",'content',['type'=>9]),6,'.','') ?: 0,
            'list'=>$data,
        );
        $this->echo_json(1,'ok',$ret);
    }

    //旷工算力
    public function minersSlAction(){
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        //用户矿池总收益
        $minerIns = $this->db->get("member_wallet",'*',['uid'=>$this->uid,'coin_id'=>$this->coinId]);
        //全网总收益
        $total = $this->db->sum("member_wallet",'miner_total',['coin_id'=>$this->coinId]);
        //持币算力  自身持币量/当天发行量
        //当天发行量
        $today = $this->db->get("coin_circulation",'*',['inputtime[>=]'=>strtotime(date('Y-m-d')),'inputtime[<=]'=>time()]);
        $hold_sl = 0;

        if($today['nums'] > 0){
            $hold_sl = $minerIns['total']/$today['nums'];
        }
        $ret = [
            'total'=>number_format($total,6,'.',''),
            'hold_sl'=>number_format($hold_sl,6,'.',''),
            'share_sl'=>number_format($minerIns['share_sl'],6,'.',''),
        ];
        $this->echo_json(1,'ok',$ret);
    }

    //旷工列表
    public function minerListAction(){
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        $page = trim($this->requestParams['page']) ?: 1;
        $list = $this->db->select("member_invite",['[>]member'=>['rid'=>'uid'],'[>]member_wallet'=>['rid'=>'uid']],[
            'member.username',
            'member_wallet.total',
            'member_wallet.share_sl(sl)',
        ],['member_wallet.coin_id'=>$this->coinId,'member_invite.uid'=>$this->uid,'LIMIT'=>[($page -1)*15,15],'ORDER'=>['member_invite.id'=>'DESC']]);
        $this->echo_json(1,'ok',$list);
    }

    //矿池总日收益
    public function minerYesterdayAction(){
        $ret = [
            'num'=>0 ,
            'coinName'=>$this->coinName,
            'rmbNum'=>0,
            'min'=>0,
            'best'=>0,
            'fee'=>0,
            'increase'=>'0',
            'netStatus'=>'1',
        ];
        if(!$this->uid){
            $this->echo_json(1,'ok',$ret);
        }
        $startTime = strtotime(date("Y-m-d 00:00:00"));
        $list = $this->db->sum("member_wallet_coin_log",'value',['uid'=>$this->uid,'type'=>107,'inputtime[>=]'=>$startTime]) ?: 0;

        //前天的收益
        $end = strtotime(date("Y-m-d 00:00:00",strtotime("-1 days")));
        $dayBeforeYes = $this->db->sum("member_wallet_coin_log",'value',['uid'=>$this->uid,'type'=>107,'inputtime[>=]'=>$end,'inputtime[<]'=>$startTime]) ?: 0;
        $increase = ($dayBeforeYes > 0) ? ($list - $dayBeforeYes)/$dayBeforeYes : 0;
        $coinInfo = $this->db->get("member_wallet_coin",'*',['id'=>$this->coinId]);
        $ret = [
            'num'=>number_format($list,6,'.',''),
            'coinName'=>$coinInfo['coinname'],
            'rmbNum'=>number_format($list*$coinInfo['usdt_rate']*$this->redis->get("huilv_USDCNY"),6,'.',''),
            'time'=>$startTime,
            'fee'=>number_format(0.1,2,'.',''),
            'min'=> $this->db->get("setting",'content',['type'=>8]) ?: 0,
            'best'=> number_format($this->db->get("setting",'content',['type'=>9]),6,'.','') ?: 0,
            'increase'=>number_format($increase,6,'.',''),
            'netStatus'=>$this->db->get("setting",'content',['type'=>10]),
        ];
        $this->echo_json(1,'ok',$ret);
    }

    //交易对添加自选
    public function collectExcpairsAction(){
        $id = intval($this->requestParams['id']); //交易对id
        if(!$id || !$this->uid){
            $this->echo_json(0,'参数错误');
        }
        $pairs = $this->db->get("exc_collect",'*',['uid'=>$this->uid,'pairs_id'=>$id]);
        if($pairs){
            //取消自选
            $res = $this->db->delete("exc_collect",['id'=>$pairs['id']]);
            if(!$res){
                $this->echo_json(0,'失败');
            }
        }else{
            //加入自选
            $res = $this->db->insert("exc_collect",array(
                'uid'=>$this->uid,
                'pairs_id'=>$id
            ));
            if(!$res){
                $this->echo_json(0,'失败');
            }
        }
        $this->echo_json(1,'ok');
    }

    //自选交易对
    public function collectListAction(){
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        $list = $this->db->select("exc_collect",'*',['uid'=>$this->uid]);
        $this->echo_json(1,'ok',$list);
    }
    //测试 pow(x,y) 函数返回 x 的 y 次方。
    public function testAction(){
	    echo strtotime(date("2020-04-23 00:00:00"));
        $username = trim($this->requestParams['username']);
        if(!preg_match("/[0-9A-Za-z]{3,16}/",$username) || strlen($username) < 3 || strlen($username) > 16){
              $this->echo_json(0,'请输入3-16位字母+数字组成的用户名');
        }
        $num = pow(9,1/2);
        echo $num;
    }

    /**
     * 内部转账
     */
    public function selfTransferAction() {
        $data = $this->requestParams['data'];
        $data['uid'] = $this->uid;
        if(empty($data['coin_id']) || empty($data['uid']) || empty($data['username'])){
            $this->echo_json(0, '缺少参数');
        }
        $jsuserinfo = $this->db->get('member','*',['OR'=>['username'=>$data['username'],'addr'=>$data['username']]]);
        if (empty($jsuserinfo)) {
            $this->echo_json(0, '对方账户不存在');
        }

        //没有实名不能转账
        $auth = $this->db->get("member_auth",'*',['uid'=>$data['uid'],'status'=>3]);
        if(empty($auth) && $data['coin_id'] == 5){
            $this->echo_json(0, '请实名后再转账');
        }

        if ($jsuserinfo['uid'] == $data['uid']) {
            $this->echo_json(0, '不需要给自己转');
        }

//        if (!$this->isp($data['uid'], $jsuserinfo['uid'])) {
//            if (!$this->isch($data['uid'], $jsuserinfo['uid'])) {
//                $this->echo_json(0, '非网体不能转账');
//            }
//        }

        $coininfo = $this->db->get('member_wallet_coin','*',['id'=>$data['coin_id']]);

        if($data['num'] <= 0){
            $this->echo_json(0, '请输入数量');
        }

//        if(empty($data['code'])){
//            $this->echo_json(0, '请输入验证码');
//        }
        if(empty($data['deal_psw'])){
            $this->echo_json(0, '请输入支付密码');
        }

//        $member = $this->dr_member_info($data['uid']);
        $member = $this->db->get("member",'*',['uid'=>$data['uid']]);
        $userinfo = $this->db->get('member_wallet_userinfo','*',['uid'=>$data['uid']]);
//        if(md5(md5($data['deal_psw']).$member['salt'].md5($data['deal_psw'])) != $userinfo['deal_psw']){
//            $this->echo_json(0, '支付密码错误');
//        }
        if(!$this->password_verify($data['deal_psw'],$userinfo['deal_psw'], $member['salt'])){
            $this->echo_json(0, '支付密码错误');
        }
        if($coininfo['isintertr'] == 2){
            //禁止提币
            $this->echo_json(0, '该币种不支持内部转账');
        }
        $wallet = $this->db->get('member_wallet','*',array('uid' => $data['uid'],'coin_id' => $data['coin_id']));
        if(empty($wallet)){
            //未开通资产
            $this->echo_json(0, '请添加资产');
        }
        //计算手续费
        $fees = 0.1;//$data['num'] * dr_block('txsxf');
        if($coininfo['id'] == 5){
            //usdt
            $fees = 0.5;
        }
        $nums = $data['num'] + $fees;
        if($wallet['total'] < $nums){
            //余额不足
            $this->echo_json(0, '余额不足');
        }
        //验证码
        $code = $this->redis->get('sms_' . $member['phone']);
        $this->redis->del('sms_' . $member['phone']);//删除验证码
        if (empty($smsuid) && !empty($data['email'])) {
            $code = $this->redis->get('email_' . $member['email']);
            $this->redis->del('email_' . $member['email']);//删除验证码
        }
//        if($code != $data['code']){
//            //验证码错误
//            $this->echo_json(0, '验证码错误');
//        }
        //删除验证码
        $this->redis->del("email_".$member['email']);
        if($this->db->action(function($db) use ($data, $nums, $wallet, $coininfo, $jsuserinfo, $fees){
            try{
                //划出余额
                $res = $db->update('member_wallet',['total[-]'=>$nums],['uid'=>$data['uid'],'coin_id'=>$data['coin_id']]);
                if(!$res){
                    throw new Exception('提币失败');
                }

                //钱包变动记录
                $this->db->insert('member_wallet_coin_log',array(
                    'uid' => $data['uid'],
                    'wallet_id' => $wallet['id'],
                    'coin_id' => $data['coin_id'],
                    'coin_name'=>$coininfo['coinname'],
                    'value' => -1 * $nums,
                    'balance' =>$wallet['total'],
                    'note' => '内部转账',
                    'type' => 2,
                    'inputtime' => time(),
                    'status' => 1,
                    'cid'=>$jsuserinfo['uid']
                ));
                $cid = $db->id();
                if($cid <= 0){
                    throw new Exception('失败');
                }

                //划进余额
                $res = $db->update('member_wallet',['total[+]'=>($nums-$fees)],array('uid' => $jsuserinfo['uid'],'coin_id' => $data['coin_id']));
                if(!$res){
                    throw new Exception('提币失败1');
                }

                //钱包变动记录
                $this->db->insert('member_wallet_coin_log',array(
                    'uid' => $jsuserinfo['uid'],
                    'wallet_id' => $wallet['id'],
                    'coin_id' => $data['coin_id'],
                    'coin_name'=>$coininfo['coinname'],
                    'value' => ($nums-$fees),
                    'balance' =>$wallet['total'],
                    'note' => '内部转账',
                    'type' => 2,
                    'inputtime' => time(),
                    'status' => 1,
                    'cid'=>$cid
                ));
                $id = $this->db->id();
                if($id <= 0){
                    throw new Exception('失败');
                }
                $res = $db->update('member_wallet_coin_log',['cid'=>$id],['id'=>$cid]);
                if(!$res){
                    throw new Exception('提币失败1');
                }
                return true;
            }catch(Exception $e){

                return false;
            }
        })){
            $this->echo_json(1, 'ok');
        }
        $this->echo_json(0, '提币失败');
    }

    //账本详情
    public function logDetailsAction(){
        $id = intval($this->requestParams['id']);
        if(!$id){
            $this->echo_json(1,'ok');
        }
        $log = $this->db->get("member_wallet_coin_log",'*',['id'=>$id]);
        $log['reUsername'] = null;
        $log['reAddr'] = null;
        $user = [];
        if(in_array($log['type'],[2])){
            if($log['value'] > 0){
                $user = $this->db->get("member_wallet_coin_log",['[>]member'=>['uid'=>'uid']],['member.username','member.addr'],['member_wallet_coin_log.id'=>$log['cid']]);
            }else{
                $user = $this->db->get("member",'*',['uid'=>$log['cid']]);
            }
        }else if(in_array($log['type'],[103,108,102])){
            $user = $this->db->get("member",'*',['uid'=>$log['cid']]);
        }else if(in_array($log['type'],[106])){
            $user = $this->db->get("member_red_bag_info",['[>]member'=>['uid'=>'uid']],['member.username','member.addr'],['member_red_bag_info.id'=>$log['id']]);
        }
        if($user){
            $log['reUsername'] = substr($user['username'],0,2).'****'.substr($user['username'],-2);
            $log['reAddr'] = substr($user['addr'],0,3).'**********'.substr($user['addr'],-3);
        }
        $this->echo_json(1,'ok',$log);
    }

    //添加关联账号
    public function addRatedAction(){
        $username = trim($this->requestParams['username']);
        $password = trim($this->requestParams['password']);
        $expiration = intval($this->requestParams['expiration']);
        if(!$this->uid || !$username || !$password || !$expiration){
            $this->echo_json(0,'参数错误');
        }
        $memberInfo = $this->db->get('member', '*', ['username' => $username]);
        if (!is_array($memberInfo) || empty($memberInfo)) {
            $this->echo_json(0, '账户不存在');
        }

        if (!$this->password_verify($password, $memberInfo['password'], $memberInfo['salt'])) {
            $this->echo_json(0, '密码不正确');
        }
        $res = $this->db->insert("member_rated",array(
            'uid'=>$this->uid,
            'cuid'=>$memberInfo['uid'],
            'expiration'=>time()+3600*24*$expiration
        ));
        if(!$res){
            $this->echo_json(0,'失败');
        }
        $this->echo_json(1,'ok');
    }

    //关联账号列表
    public function ratedListAction(){
        $page = intval($this->requestParams['page']) ?: 1;
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        $list = $this->db->select("member_rated",['[>]member'=>['cuid'=>'uid']],['member.username','member.avatar','member.uid(cuid)'],['member_rated.uid'=>$this->uid,'member_rated.expiration[>]'=>time(),'LIMIT'=>[($page-1)*10,10]]);
        //
        $list1 = $this->db->select("member_rated",['[>]member'=>['uid'=>'uid']],['member.username','member.avatar','member.uid(cuid)'],['member_rated.cuid'=>$this->uid,'member_rated.expiration[>]'=>time(),'LIMIT'=>[($page-1)*10,10]]);
        if($page == 1 && $list1){
            $list = array_merge($list1,$list);
        }
        $main = $this->db->get("member",['username','avatar','uid'],['uid'=>$this->uid]);
        $res = [
            'main'=>$main,
            'list'=>$list
        ];
        $this->echo_json(1,'ok',$res);
    }

    //切换账号
    public function switchAccountAction(){
        $cuid = intval($this->requestParams['cuid']);  //需要登录的账号的uid
        $password = trim($this->requestParams['password']);  //登录密码
        if(!$this->uid || !$cuid || !$password){
            $this->echo_json(0,'参数失败');
        }
        $user = $this->db->get("member",'*',['uid'=>$this->uid]);
        $memberInfo = $this->db->get("member",'*',['uid'=>$cuid]);
        if(empty($memberInfo) || empty($user)){
            $this->echo_json(0,'用户不存在');
        }
        if (!$this->password_verify($password, $memberInfo['password'], $memberInfo['salt'])) {
            $this->echo_json(0, '密码不正确');
        }
        unset($memberInfo['password']);
        unset($memberInfo['token']);
        $char = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789';
        $token = '';
        for($i = 0; $i < 32; $i++) {
            $token .= $char[mt_rand(0, 53)];
        }
        $memberInfo['token'] = $token;
        $this->redis->set($memberInfo['token'], $memberInfo['uid'],1800);
        //删除现在账号的登录
        $this->redis->del($user['token']);
        $this->db->update('member', ['token' => $memberInfo['token']], ['uid' => $memberInfo['uid']]);
        $yaoqm = $this->db->get('member_data',['yaoqingma','yaoqingmaright','area'],['uid'=>$memberInfo['uid']]);
        $this->echo_json(1, 'ok',array_merge($memberInfo, $yaoqm));
    }

    //修改手机号
    public function changePhoneAction(){
        $deal_pwd = trim($this->requestParams['deal_pwd']);
        $idCard = trim($this->requestParams['idCard']);
        $oldCode = trim($this->requestParams['oldCode']);
        $newPhone = trim($this->requestParams['newPhone']);
        $newCode = trim($this->requestParams['newCode']);
        if(!$this->uid || !$deal_pwd || !$oldCode || !$newPhone || !$newCode){
            $this->echo_json(0,'请填写完整信息');
        }
        //是否实名
        $auth = $this->db->get("member_auth",'*',['uid'=>$this->uid,'status'=>3]);
        if($auth){
            if(!$idCard){
                $this->echo_json(0,'请填写身份证号');
            }
            if($idCard != $auth['auth_sn']){
                $this->echo_json(0,'身份证号错误');
            }
        }
        $memberInfo = $this->db->get("member",'*',['uid'=>$this->uid]);
        $userinfo = $this->db->get('member_wallet_userinfo','*',['uid'=>$this->uid]);
        if(!$this->password_verify($deal_pwd,$userinfo['deal_psw'], $memberInfo['salt'])){
            $this->echo_json(0, '支付密码错误');
        }
        //验证旧手机验证码
        $oldCodeV = $this->redis->get("sms_".$memberInfo['phone']);
        if($oldCodeV != $oldCode){
            $this->echo_json(0, '旧手机验证码错误');
        }
        $this->redis->del("sms_".$memberInfo['phone']);
        //验证新手机验证码
        $newCodeV = $this->redis->get("sms_".$newPhone);
        if($newCodeV != $newCode){
            $this->echo_json(0, '新手机验证码错误');
        }
        $this->redis->del("sms_".$newPhone);
        //修改手机号
        $res = $this->db->update("member",['phone'=>$newPhone],['uid'=>$memberInfo['uid']]);
        if(!$res){
            $this->echo_json(0, '修改手机号失败');
        }
        unset($memberInfo['password']);
        $memberInfo['phone'] = $newPhone;
        $yaoqm = $this->db->get('member_data',['yaoqingma','yaoqingmaright','area'],['uid'=>$memberInfo['uid']]);
        $this->echo_json(1, '修改手机号成功',array_merge($memberInfo, $yaoqm));
    }

    //根据用户地址获取用户平台币的交易明细
    public function userCapitalFlowAction(){
        $addr = trim($this->requestParams['addr']);
        $page = intval($this->requestParams['page']) ?: 1;
        $balance = 0;
        if($addr){
            $user = $this->db->get("member",'*',['OR'=>['addr'=>$addr,'username'=>$addr]]);
            if($user){
                $where['uid'] = $user['uid'];
                //余额
                $balance = $this->db->get("member_wallet",'total',['uid'=>$user['uid'],'coin_id'=>$this->coinId]);
            }
        }
        $where['coin_id'] = $this->coinId;
        $where['ORDER'] = ['id'=>'DESC'];
        $where['LIMIT'] = [($page -1)*10,10];
        $list = $this->db->select("member_wallet_coin_log",'*',$where);
        foreach ($list as $k=>$v) {
            $list[$k]['addr'] = $this->db->get("member",'addr',['uid'=>$v['uid']]);
            $list[$k]['txCount'] = substr(floor($v['id']/$this->coinId),0,1);
        }
        $ret = [
            'balance'=> $balance ?: 0,
            'list'=> $list,
        ];
        $this->echo_json(1,'ok',$ret);
    }

    //激活的金额  币种
    public function activeAmountAction(){
        $this->echo_json(1,'ok',['num'=>$this->activeFee,'coinName'=>$this->coinName,'send'=>$this->activeSend]);
    }

    //红包币种余额等信息
    public function redBagCoinAction(){
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        $wallet = $this->db->get("member_wallet",'*',['uid'=>$this->uid,'coin_id'=>$this->coinId]);
        $wallet['coinName'] = $this->coinName;
        $this->echo_json(1,'ok',$wallet);
    }

    //红包记录
    public function redBagLogAction(){
        $page = intval($this->requestParams['page']) ?: 1;
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        $list = $this->db->select("member_red_bag",['id','uid','amount','inputtime','type'],['LIMIT'=>[($page - 1)*10,10],'uid'=>$this->uid,'ORDER'=>['id'=>'DESC']]);
        $noteType = [
            '1'=>'邀请注册红包',
            '2'=>'主题红包',
        ];
        foreach ($list as $k=>$v){
            $list[$k]['note'] = $noteType[$v['type']];
            $list[$k]['coinName'] = $this->coinName;
        }
        $this->echo_json(1,'ok',$list);
    }

    //usdt兑换gsv
    public function transferAction(){
        if(intval($this->db->get("setting",'content',['type'=>10])) != 1){
            $this->echo_json(0,'认购活动已关闭');
        }
        $num = floatval($this->requestParams['num']);
        $deal_pwd = trim($this->requestParams['deal_pwd']);
        if(!$this->uid){
            $this->echo_json(0,'参数错误');
        }
        if($num <= 0){
            $this->echo_json(0,'请输入有效的兑换额度');
        }
        if(empty($deal_pwd)){
            $this->echo_json(0,'请输入支付密码');
        }
        $member = $this->db->get("member",'*',['uid'=>$this->uid]);
        $userinfo = $this->db->get('member_wallet_userinfo','*',['uid'=>$this->uid]);
        if(!$this->password_verify($deal_pwd,$userinfo['deal_psw'], $member['salt'])){
            $this->echo_json(0, '支付密码错误');
        }
        //可兑换额度
        $getCoinInfo = $this->db->get("member_wallet_coin",'*',['id'=>6]);
        $payCoinInfo = $this->db->get("member_wallet_coin",'*',['id'=>5]);
        $getCoinWallet = $this->db->get("member_wallet",'*',['uid'=>$this->uid,'coin_id'=>6]);
        $payCoinWallet = $this->db->get("member_wallet",'*',['uid'=>$this->uid,'coin_id'=>5]);
        if($getCoinWallet['convert_amount'] < $num){
            $this->echo_json(0,'可兑换额度不足');
        }
        $payNum = $num*$getCoinInfo['usdt_rate'];
        if($payNum <= 0 || $payCoinWallet['total'] < $payNum){
            $this->echo_json(0,'支付余额不足');
        }
        $member = $this->db->get("member",'*',['accountType'=>2]);
        if($member){
            $walletTotal = $this->db->get("member_wallet",'*',['uid'=>$member['uid'],'coin_id'=>$this->coinId]);
        }
        if(empty($walletTotal) || $walletTotal['total'] < $num){
            $this->echo_json(0,'系统可兑换额度不足');
        }
        $res = $this->db->action(function () use($getCoinWallet,$payCoinWallet,$getCoinInfo,$payCoinInfo,$num,$payNum,$walletTotal){
            try{

                //主账号减推广额度
                $update['total[-]'] = $num;
                $res = $this->db->update("member_wallet",$update,['id'=>$walletTotal['id']]);
                if(!$res){
                    throw new Exception("主账户减金额失败");
                }
                $walletTotal['total'] -=  $num;
                $res = $this->db->insert("member_wallet_coin_log",array(
                    'uid'=>$walletTotal['uid'],
                    'coin_id'=>$this->coinId,
                    'coin_name'=>$this->coinName,
                    'type'=>112,
                    'note'=>'认购兑换',
                    'value'=> -1 * $num,
                    'balance'=> $walletTotal['total'],
                    'wallet_id'=> $walletTotal['id'],
                    'status'=>1,
                    'inputtime'=>time(),
                ));
                if(!$res){
                    throw new Exception("加资金记录失败12");
                }

                //加兑换额度
                $res = $this->db->update("member_wallet",['converted_amount[+]'=>$num,'convert_amount[-]'=>$num],['id'=>$getCoinWallet['id']]);
                if(!$res){
                    throw new Exception("加兑换数量失败");
                }
                //资金记录
                $res = $this->db->insert("member_wallet_coin_log",array(
                    'uid'=>$this->uid,
                    'coin_id'=>$getCoinInfo['id'],
                    'coin_name'=>$getCoinInfo['coinname'],
                    'type'=>112,
                    'note'=>$payNum.' '.$payCoinInfo['coinname'].'认购'.$num.' '.$getCoinInfo['coinname'],
                    'value'=> $num,
                    'balance'=>$getCoinWallet['total'],
                    'wallet_id'=>$getCoinWallet['id'],
                    'status'=>1,
                    'inputtime'=>time(),
                ));
                if(!$res){
                    throw new Exception("加资金记录失败");
                }
                //减金额
                $res = $this->db->update("member_wallet",['total[-]'=>$payNum],['id'=>$payCoinWallet['id']]);
                if(!$res){
                    throw new Exception("加兑换数量失败");
                }
                //资金记录
                $res = $this->db->insert("member_wallet_coin_log",array(
                    'uid'=>$this->uid,
                    'coin_id'=>$payCoinInfo['id'],
                    'coin_name'=>$payCoinInfo['coinname'],
                    'type'=>112,
                    'note'=>$payNum.' '.$payCoinInfo['coinname'].'认购'.$num.' '.$getCoinInfo['coinname'],
                    'value'=> -1 * $payNum,
                    'balance'=>$payCoinWallet['total'],
                    'wallet_id'=>$payCoinWallet['id'],
                    'status'=>1,
                    'inputtime'=>time(),
                ));
                if(!$res){
                    throw new Exception("加资金记录失败");
                }
                return true;
            }catch (Exception $e){
                return false;
            }
        });
        if($res === true){
            $this->echo_json(1,'成功');
        }
        $this->echo_json(0,'失败');
    }

    //获取当前用户待处理的发布和订单
    public function waitingDealAction(){
        if(!$this->uid){
            $this->echo_json(1,'ok');
        }
        $listNums = 0;
        $list = $this->db->select("trans_list",'*',['uid'=>$this->uid,'status'=>[1,2,3,4]]);
        foreach ($list as $v){
            $listNums += $this->db->count("trans_order",['cid'=>$v['id'],'status'=>[1,2,3]]);
        }
        $order = $this->db->count("trans_order",['uid'=>$this->uid,'status'=>[1,2,3]]);
        $this->echo_json(1,'ok',['order'=>$order?:0,'list'=>$listNums]);
    }

    //系统某一交易对所有的委托记录
    public function tradeLogAction(){
        $status = trim($this->requestParams['status']);
        $page = intval($this->requestParams['page']) ?: 1;
        $excpairs = intval($this->requestParams['excpairs']);
        $where = [];
        if($excpairs){
            $where['excpairs'] = $excpairs;
        }
        if($status){
            $where['status'] = explode(',',$status);
        }
        $where['LIMIT'] = [($page - 1)*2,2];
        $where['ORDER'] = ['id'=>'DESC'];
        $list = $this->db->select("exc_entrust_log",'*',$where);
        foreach($list as $k=>$v){
            $list[$k]['username'] = $this->db->get("member",'username',['uid'=>$v['uid']]);
            $list[$k]['inputtime'] = date("Y-m-d H:i:s");
        }
        $this->echo_json(1,'ok',$list);
    }

    //算力数据
    public function slAction(){
        //两个总账号余额
        $memberUid = $this->db->select("member",'*',['accountType[!]'=>1]);
        $wallet = [];
        if($memberUid){
            $wallet = $this->db->select("member_wallet",'*',['uid'=>$memberUid,'coin_id'=>6]);
        }
        //本台GSV持有者
        $hold = $this->db->count("member_wallet",'*',['coin_id'=>6,'total[>]'=>0]);
        //提现交易数量
        $cashNums = 0;
//        $cashNums = $this->db->sum("member_wallet_cash_log",'value',['status'=>2]);

        //系统内部转账的笔数  内部转账  提现-内部转账
        $transCount  = $this->db->count("member_wallet_coin_log",['OR'=>['note'=>['内部转账','提现-内部转账'],'AND'=>['type'=>107,'value[>]'=>0]]]);
        //充值交易数量
        $rechargeNums = 0;
//        $rechargeNums = $this->db->sum("member_wallet_translate_log",'amount');

        //当日产出
        $today = $this->db->get("coin_circulation",'*',['inputtime[<]'=>strtotime(date("Y-m-d 00:00:00",strtotime('+1 days'))),'inputtime[>=]'=>strtotime(date('Y-m-d 00:00:00'))]);
        $res = [
            'wallet'=>$wallet,
            'hold'=>$hold,
            'cashNums'=> abs($cashNums) ?: 0,
            'transCount'=> $transCount ?: 0,
            'rechargeNums' => $rechargeNums ,
            'today' => $today['nums'] ?: 0 ,
        ];
        $this->echo_json(1,'ok',$res);
    }

    //最新交易（内部转账数据）
    public function latestTransAction(){
        $list = $this->db->select("member_wallet_coin_log",'*',['type'=>2,'note'=>'内部转账','value[>]'=>0,'ORDER'=>['id'=>'DESC'],'LIMIT'=>4]);
        foreach ($list as $k=>$v){
            $list[$k]['addr'] = $this->db->get("member",'addr',['uid'=>$v['uid']]);
        }
        $this->echo_json(1,'ok',$list);
    }

    //-------------------------------脚本----------------------------//
    //计算推广算力
    public $total = 0;
    public function calShareSlAction(){
        $page = 1;
        while (true){
            $member = $this->db->select("member",'uid',['accountType[!]'=>1]);
            if($member){
                $where['uid[!]'] = $member;
            }
            $where['LIMIT'] = [($page -1)*10,10];
//            $where['uid'] = [2114,2115,2116,2117,2118,2119,2120,2113];
            $list = $this->db->select("member",'*',$where);
            if(empty($list)){
               sleep(3);
               $page = 1;
//               break;
               continue;
            }
            foreach($list as $v){
                $invite = $this->db->select("member_invite",'*',['uid'=>$v['uid']]);
                $data = [];
                foreach($invite as $val){
                    $this->total = 0;
                    $this->userSl($val['rid']);
                    $data[] = $this->total;
                }
                if(empty($data)){
                    continue;
                }
                rsort($data);
                $shareSl =  pow($data[0],1/3);
                for($i = 1;$i < count($data);$i++){
                    if($data[$i] > 10000){
                        $shareSl += 10000*10+($data[$i] - 10000);
                    }else{
                        $shareSl += $data[$i]*10;
                    }
                }
                if($shareSl > 0){
                    //修改用户推广算力
                    $this->db->update("member_wallet",['share_sl'=>$shareSl],['uid'=>$v['uid'],'coin_id'=>$this->coinId]);
                }
            };
            $page++;
        }
    }

    public function userSl($uid){
        $wallet = $this->db->get("member_wallet",'*',['uid'=>$uid,'coin_id'=>$this->coinId]);
        $this->total += $wallet['total'];
        $list = $this->db->select("member_invite",'*',['uid'=>$uid]);
        foreach($list as $v){
            $this->userSl($v['rid']);
        }
    }

    //记录每天的余额
    public function balanceLogAction(){
        $status = intval($this->db->get("setting",'content',['type'=>14]));
        if($status != 1){
            echo "收益未开启";exit;
        }
        //最小持币量
        $min = $this->db->get("setting",'content',['type'=>8]) ?: 0;
        $page = 1;
        $last = [];
        $num = 0;
        $member = $this->db->select("member",'uid',['accountType[!]'=>1]);
        while (true){
            $where['coin_id'] = $this->coinId;
            $where['ORDER'] = ['total'=>'ASC','id'=>'ASC'];
            $where['total[>=]'] = $min;
            $where['LIMIT'] = [($page - 1)*50,50];
            if($member){
                $where['uid[!]'] = $member;
            }
            $list  = $this->db->select("member_wallet",'*',$where);
            if(empty($list)){
                break;
            }
            for($i = 0;$i < count($list);$i++){
                $num++;
                if($num == 1 && empty($last)){
                    $list[$i]['rank'] = 1;
                }else{
                    if($list[$i]['total'] > $last['total']){
//                    $list[$i]['rank'] = $list[$i-1]['rank']+1;
                        $list[$i]['rank'] = $num;
                    }else{
                        $list[$i]['rank'] = $last['rank'];
                    }
                }
                $v = $list[$i];
                $v['wallet_id'] = $v['id'];
                unset($v['updatetime']);
                unset($v['id']);
                $v['in_status'] = 1;
                $v['inputtime'] = time();
                $v['zerotime'] = strtotime(date("Y-m-d 00:00:00"));
                $res = $this->db->insert("member_wallet_balance",$v);
                if(!$res){
                    $content = date("Y-m-d H:i:s").', 失败， 记录余额：'.json_encode($v).'， 金额：0'.json_encode($this->db->error());
                    file_put_contents(dirname(__FILE__).'/miner_rank_error'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                }else{
                    $content = date("Y-m-d H:i:s").', 成功， 记录余额：'.json_encode($v).'， 金额：0';
                    file_put_contents(dirname(__FILE__).'/miner_rank'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                }
                $last = $list[$i];
            }
            $page++;
        }
        //算力
        $page1 = 1;
        while (true){
            $where1['coin_id'] = $this->coinId;
            $where1['ORDER'] = ['id'=>'ASC'];
            $where1['total[<]'] = $min;
            $where1['share_sl[>]'] = 0;
            $where1['LIMIT'] = [($page1 - 1)*50,50];
            if($member){
                $where1['uid[!]'] = $member;
            }
            $list1 = $this->db->select("member_wallet",'*',$where1);
            if(empty($list1)){
                break;
            }
            foreach($list1 as $k=>$v){
                $v['wallet_id'] = $v['id'];
                unset($v['updatetime']);
                unset($v['id']);
                $v['inputtime'] = time();
                $v['rank'] = 0;
                $v['in_status'] = 1;
                $v['zerotime'] = strtotime(date("Y-m-d 00:00:00"));
                $res = $this->db->insert("member_wallet_balance",$v);
                if(!$res){
                    $content = date("Y-m-d H:i:s").', 失败， 记录余额：'.json_encode($v).'， 金额：0'.json_encode($this->db->error());
                    file_put_contents(dirname(__FILE__).'/miner_rank_error'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                    continue;
                }
                $content = date("Y-m-d H:i:s").', 成功， 记录余额：'.json_encode($v).'， 金额：0';
                file_put_contents(dirname(__FILE__).'/miner_rank'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
            }
            $page1++;
        }
    }
    //矿池持币收益
    public function rankAction(){
        $status = intval($this->db->get("setting",'content',['type'=>14]));
        if($status != 1){
            echo "收益未开启";exit;
        }
        $member = $this->db->get("member",'*',['accountType'=>3]);
        if(empty($member)){
            echo "主账户错误";exit;
        }
        $walletTotal = $this->db->get("member_wallet",'*',['coin_id'=>$this->coinId,'uid'=>$member['uid']]);
        if($walletTotal['total'] <= 0){
            echo "主账户余额不足";exit;
        }
        $where['coin_id'] = $this->coinId;
        $where['in_status'] = 1;
//        $where['inputtime[<=]'] = strtotime(date('Y-m-d 00:00:00'));
//        $where['inputtime[>=]'] = strtotime(date("Y-m-d 00:00:00",strtotime('-1 days')));
        $ranks = $this->db->sum("member_wallet_balance",'rank',$where);
        $slTotal = $this->db->sum("member_wallet_balance",'share_sl',$where);
//        //总排名
//        $ranks = array_sum(array_column($list,'rank'));
//        //全网总算力
//        $slTotal = array_sum(array_column($list,'share_sl'));
        $today = $this->db->get("coin_circulation",'*',['inputtime[>=]'=>strtotime(date("Y-m-d 00:00:00",strtotime('-1 days'))),'inputtime[<]'=>strtotime(date('Y-m-d 00:00:00'))]);
        if(empty($today)){
            echo "没有发行量";exit;
        }
        $page = 1;
        $bestCoin = 0;
        $bestCoinIns = 0;
        while (true){
//            $where['ORDER'] = ['rank'=>'DESC'];
            $where['LIMIT'] = [($page -1)*10,10];
            $list = $this->db->select("member_wallet_balance",'*',$where);
            if(empty($list)){
                break;
            }
            foreach($list as $k=>$v){
                if($walletTotal['total'] <= 0){
                    $content = date("Y-m-d H:i:s").', 矿池持币收益主账号金额不足';
                    file_put_contents(dirname(__FILE__).'/miner_rank.txt', $content.PHP_EOL,FILE_APPEND);
                    break;
                }
                //修改收益发放状态
                $res = $this->db->update("member_wallet_balance",['in_status'=>2],['id'=>$v['id']]);
                if(!$res){
                    $content = date("Y-m-d H:i:s").', 修改收益发放状态失败:'.json_encode($v);
                    file_put_contents(dirname(__FILE__).'/miner_rank.txt', $content.PHP_EOL,FILE_APPEND);
                    continue;
                }

                $ins = ($v['rank']/$ranks)*($today['nums']/2);
                if($ins > 0 && $v['rank'] > 0){
                    if($ins > $bestCoinIns){
                        $bestCoinIns = $ins;
                        $bestCoin = $v['total'];
                    }
                    $res = $this->db->action(function() use($ins,$v,$walletTotal){
                        try{

                            $update = [];
                            //主账户减金额
                            $update['total[-]'] = $ins;
                            if($walletTotal['total'] < $ins){
                                $update['total[-]'] = $walletTotal['total'];
                            }
                            $res = $this->db->update("member_wallet",$update,['id'=>$walletTotal['id']]);
                            if(!$res){
                                throw new Exception("主账户减金额失败");
                            }
                            $walletTotal['total'] -=  $update['total[-]'];
                            $res = $this->db->insert("member_wallet_coin_log",array(
                                'uid'=>$walletTotal['uid'],
                                'coin_id'=>$this->coinId,
                                'coin_name'=>$this->coinName,
                                'type'=> 107,
                                'note'=> '矿池持币收益',
                                'value'=> -1 * $update['total[-]'],
                                'balance'=> $walletTotal['total'],
                                'wallet_id'=> $walletTotal['id'],
                                'status'=>1,
                                'inputtime'=>time(),
                            ));
                            if(!$res){
                                throw new Exception("加资金记录失败12");
                            }


                            //添加收益
                            $res = $this->db->update("member_wallet",['total[+]'=>$ins,'miner_total[+]'=>$ins],['id'=>$v['wallet_id']]);
                            if(!$res){
                                throw new Exception("加资金记录失败4");
                            }
                            $res = $this->db->insert("member_wallet_coin_log",array(
                                'uid'=>$v['uid'],
                                'coin_id'=>$this->coinId,
                                'coin_name'=>$this->coinName,
                                'type'=>107,
                                'note'=>'矿池持币收益',
                                'value'=> $ins,
                                'balance'=>$v['total'],
                                'wallet_id'=>$v['wallet_id'],
                                'status'=>1,
                                'inputtime'=>time()
                            ));
                            if(!$res){
                                throw new Exception("加资金记录失败3");
                            }
                            return true;
                        }catch (Exception $e){
                            return false;
                        }
                    });
                    if($res === true){
                        $content = date("Y-m-d H:i:s").', 成功， 矿池持币收益：'.json_encode($v).'， 金额：'.$ins;
                        file_put_contents(dirname(__FILE__).'/miner_rank_ins'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                    }else{
                        $content =  $content = date("Y-m-d H:i:s").'，失败， 矿池持币收益：'.json_encode($v).'， 金额：'.$ins;
                        file_put_contents(dirname(__FILE__).'/miner_rank_ins_error'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                    }
                }

                //推广收益
                if($v['share_sl'] > 0 && $slTotal > 0 && $walletTotal['total'] > 0){
                    $shareIns =  ($v['share_sl']/$slTotal)*($today['nums']/2);
                    if($shareIns > 0){
                        $res = $this->db->action(function() use($shareIns,$v,$ins,$walletTotal){
                            try{

                                //主账户减金额
                                $update = [];
                                $update['total[-]'] = $shareIns;
                                if($walletTotal['total'] < $shareIns){
                                    $update['total[-]'] = $walletTotal['total'];
                                }
                                $res = $this->db->update("member_wallet",$update,['id'=>$walletTotal['id']]);
                                if(!$res){
                                    throw new Exception("主账户减金额失败");
                                }
                                $walletTotal['total'] -=  $update['total[-]'];
                                $res = $this->db->insert("member_wallet_coin_log",array(
                                    'uid'=>$walletTotal['uid'],
                                    'coin_id'=>$this->coinId,
                                    'coin_name'=>$this->coinName,
                                    'type'=> 107,
                                    'note'=> '矿池推广收益',
                                    'value'=> -1 * $update['total[-]'],
                                    'balance'=> $walletTotal['total'],
                                    'wallet_id'=> $walletTotal['id'],
                                    'status'=>1,
                                    'inputtime'=>time(),
                                ));
                                if(!$res){
                                    throw new Exception("加资金记录失败12");
                                }


                                //添加收益
                                $res = $this->db->update("member_wallet",['total[+]'=>$shareIns,'miner_total[+]'=>$shareIns],['id'=>$v['wallet_id']]);
                                if(!$res){
                                    throw new Exception("加资金记录失败");
                                }
                                $res = $this->db->insert("member_wallet_coin_log",array(
                                    'uid'=>$v['uid'],
                                    'coin_id'=>$this->coinId,
                                    'coin_name'=>$this->coinName,
                                    'type'=>107,
                                    'note'=>'矿池推广收益',
                                    'value'=> $shareIns,
                                    'balance'=> ($v['total']+$ins),
                                    'wallet_id'=>$v['wallet_id'],
                                    'status'=>1,
                                    'inputtime'=>time()
                                ));
                                if(!$res){
                                    throw new Exception("加资金记录失败");
                                }
                                return true;
                            }catch (Exception $e){
                                return false;
                            }
                        });
                        if($res === true){
                            $content = date("Y-m-d H:i:s").', 成功， 矿池推广收益：'.json_encode($v).'， 金额：'.$shareIns;
                            file_put_contents(dirname(__FILE__).'/miner_rank_ins'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                        }else{
                            $content =  $content = date("Y-m-d H:i:s").'，失败， 矿池推广收益：'.json_encode($v).'， 金额：'.$shareIns;
                            file_put_contents(dirname(__FILE__).'/miner_rank_ins_error'.date('Y-m-d').'.txt', $content.PHP_EOL,FILE_APPEND);
                        }
                    }
                }
                if($today['best'] && $today['best'] > 0){
                    $bestCoin = $today['best'];
                }
                if($bestCoin > 0){
                    $this->db->update("setting",['content'=>$bestCoin],['type'=>9]);
                }
            }
//            $page++;
        }
    }
    //过期红包退回
    public function refundAction(){
        $page = 1;
        while (true){
            $where['LIMIT'] = [($page -1)*10,10];
            $where['expiration[<=]'] = time();
            $where['left_num[>]'] = 0;
            $where['status'] = 1;
            $list = $this->db->select("member_red_bag",'*',$where);
            if(empty($list)){
                sleep(3);
                $page = 1;
                continue;
            }
            foreach($list as $k=>$v){
               $this->db->action(function() use ($v){
                    try{
                        //获取当前红包没有领取的金额
                        $data = $this->db->select("member_red_bag_info",'*',['bag_id'=>$v['id'],'status'=>1]);
                        //退回金额
                        $amount = array_sum(array_column($data,'amount'));
                        $ids = array_column($data,'id');
                        $wallet = $this->db->get("member_wallet",'*',['uid'=>$v['uid'],'coin_id'=>$this->coinId]);
                        if(empty($wallet)){
                            throw new Exception("资金账户不存在");
                        }
                        if($amount > 0){
                            $res = $this->db->update("member_wallet",['total[+]'=>$amount],['id'=>$wallet['id']]);
                            if(!$res){
                                throw new Exception("加金额失败");
                            }
                            $res = $this->db->insert("member_wallet_coin_log",array(
                                'uid'=>$v['uid'],
                                'coin_id'=>$this->coinId,
                                'coin_name'=>$this->coinName,
                                'type'=>109,
                                'note'=>'红包过期退回',
                                'value'=> $amount,
                                'balance'=>$wallet['total'],
                                'wallet_id'=>$wallet['id'],
                                'status'=>1,
                                'inputtime'=>time(),
                                'cid'=>$v['id']
                            ));
                            if(!$res){
                                throw new Exception("加资金记录失败");
                            }
                        }
                        //修改红包状态
                        $res = $this->db->update("member_red_bag_info",['status'=>2],['id'=>$ids]);
                        if($res != count($ids)){
                            throw new Exception("修改红包状态失败");
                        }
                        $res = $this->db->update("member_red_bag",['status'=>2],['id'=>$v['id']]);
                        if(!$res){
                            throw new Exception("修改红包状态失败1");
                        }
                        $content = date("Y-m-d H:i:s").', 成功， 红包退回：'.json_encode($v).'， 金额：'.$amount;
                        file_put_contents(dirname(__FILE__).'/reg_bag_refund_'.date("Y-m-d").'.txt', $content.PHP_EOL,FILE_APPEND);
                        return true;
                    }catch (Exception $e){
                        $content = date("Y-m-d H:i:s").', 失败， 红包退回：'.json_encode($v).'， 错误信息：'.$e->getMessage();
                        file_put_contents(dirname(__FILE__).'/reg_bag_refund_'.date("Y-m-d").'.txt', $content.PHP_EOL,FILE_APPEND);
                        return false;
                    }
                });
            }
        }
    }

    //激活后赠送100U的兑换额度
    public function activeSendAction(){
        $page = 1;
        $id = 0;
        while (true){
            $coinInfo = $this->db->get("member_wallet_coin",'*',['id'=>$this->coinId]);
            if($coinInfo['usdt_rate'] <= 0){
                continue;
            }
            $where['LIMIT'] = [($page - 1)*10,10];
            $where['id[>=]'] = $id;
            $where['type'] = 108;
            $list = $this->db->select("member_wallet_coin_log",'*',$where);
            if(empty($list)){
                sleep(3);
                $page = 1;
                continue;
            }
            foreach($list as $k=>$v){
                $id = $v['id'];
                //当前账号是否已经赠送
                $log = $this->db->get("member_wallet_coin_log",'*',['cid'=>$v['id'],'type'=>110,'uid'=>$v['uid']]);
                if($log){
                    continue;
                }
                $this->db->action(function() use($v,$coinInfo){
                    $num = 100/$coinInfo['usdt_rate'];
                    try{
                        //加兑换额度
                        $wallet = $this->db->get("member_wallet",'*',['coin_id'=>$coinInfo['id'],'uid'=>$v['uid']]);
                        if(empty($wallet)){
                            $res = $this->db->insert("member_wallet",array(
                                'uid'=>$v['uid'],
                                'coin_id'=>$coinInfo['id'],
                                'total'=>0,
                                'frozen'=>0
                            ));
                            if(!$res){
                                throw new Exception("加资金账户失败");
                            }
                            $wallet['id'] = $this->db->id();
                        }
                        $res = $this->db->update("member_wallet",['convert_amount[+]'=>$num],['id'=>$wallet['id']]);
                        if(!$res){
                            throw new Exception("加兑换额度失败");
                        }
                        //资金记录
                        $res = $this->db->insert("member_wallet_coin_log",array(
                            'uid'=>$v['uid'],
                            'coin_id'=>$coinInfo['id'],
                            'coin_name'=>$coinInfo['coinname'],
                            'type'=>110,
                            'note'=>'激活账号送认购额度',
                            'value'=> $num,
                            'balance'=>$wallet['total'],
                            'wallet_id'=>$wallet['id'],
                            'status'=>1,
                            'inputtime'=>time(),
                            'cid'=>$v['id']
                        ));
                        if(!$res){
                            throw new Exception("加资金记录失败");
                        }
                        $content = date("Y-m-d H:i:s").', 成功， 激活送认购额度：'.json_encode($v).'， 金额：'.$num;
                        file_put_contents(dirname(__FILE__).'/miner_rank.txt', $content.PHP_EOL,FILE_APPEND);
                        return true;
                    }catch (Exception $e){
                        $content = date("Y-m-d H:i:s").', 失败，错误信息：'.$e->getMessage().' 激活送认购额度：'.json_encode($v).'， 金额：'.$num;
                        file_put_contents(dirname(__FILE__).'/miner_rank.txt', $content.PHP_EOL,FILE_APPEND);
                        return false;
                    }
                });
            }
            $page++;
        }
    }

    //下级实名认证后上级赠送15U 认购额度
    public function authSendAction(){
        if(intval($this->db->get("setting",'content',['type'=>11])) != 1){
            echo "实名赠送额度关闭";exit;
        }
        $page = 1;
        while(true){
            $coinInfo = $this->db->get("member_wallet_coin",'*',['id'=>$this->coinId]);
            if($coinInfo['usdt_rate'] <= 0){
                continue;
            }
            $where['LIMIT'] = [($page - 1)*10,10];
            $where['upStatus'] = 1;
            $where['status'] = 3;
            $list = $this->db->select("member_auth",'*',$where);
            if(empty($list)){
                sleep(3);
                continue;
            }
            foreach ($list as $k=>$v){
                $this->db->action(function() use($v,$coinInfo){
                    try{
                        if(time() >= strtotime(date("2020-04-25 00:00:00"))){
                            $num = 5/$coinInfo['usdt_rate'];
                        }else{
                            $num = 7.5/$coinInfo['usdt_rate'];
                        }
                        //上级
                        $sup = $this->db->get("member_invite",'*',['rid'=>$v['uid']]);
                        if($sup){
                            //是否已经领取
                            $log = $this->db->get("member_wallet_coin_log",'*',['cid'=>$v['uid'],'type'=>111]);
                            if(!$log){

                                //加兑换额度
                                $wallet = $this->db->get("member_wallet",'*',['coin_id'=>$coinInfo['id'],'uid'=>$sup['uid']]);
                                if(empty($wallet)){
                                    $res = $this->db->insert("member_wallet",array(
                                        'uid'=>$sup['uid'],
                                        'coin_id'=>$coinInfo['id'],
                                        'total'=>0,
                                        'frozen'=>0
                                    ));
                                    if(!$res){
                                        throw new Exception("加资金账户失败");
                                    }
                                    $wallet['id'] = $this->db->id();
                                }
                                $res = $this->db->update("member_wallet",['convert_amount[+]'=>$num],['id'=>$wallet['id']]);
                                if(!$res){
                                    throw new Exception("加兑换额度失败");
                                }
                                //资金记录
                                $res = $this->db->insert("member_wallet_coin_log",array(
                                    'uid'=>$sup['uid'],
                                    'coin_id'=>$coinInfo['id'],
                                    'coin_name'=>$coinInfo['coinname'],
                                    'type'=>110,
                                    'note'=>'下级实名认证送认购额度',
                                    'value'=> $num,
                                    'balance'=>$wallet['total'],
                                    'wallet_id'=>$wallet['id'],
                                    'status'=>1,
                                    'inputtime'=>time(),
                                    'cid'=>$v['uid']
                                ));
                                if(!$res){
                                    throw new Exception("加资金记录失败");
                                }
                            }
                        }

                        if(time() >= strtotime(date("2020-04-25 00:00:00"))){
                            $num1 = 50/$coinInfo['usdt_rate'];
                        }else{
                            $num1 = 75/$coinInfo['usdt_rate'];
                        }
                        //增加自身额度
                        //加兑换额度
                        $wallet = $this->db->get("member_wallet",'*',['coin_id'=>$coinInfo['id'],'uid'=>$v['uid']]);
                        if(empty($wallet)){
                            $res = $this->db->insert("member_wallet",array(
                                'uid'=>$v['uid'],
                                'coin_id'=>$coinInfo['id'],
                                'total'=>0,
                                'frozen'=>0
                            ));
                            if(!$res){
                                throw new Exception("加资金账户失败");
                            }
                            $wallet['id'] = $this->db->id();
                        }
                        $res = $this->db->update("member_wallet",['convert_amount[+]'=>$num1],['id'=>$wallet['id']]);
                        if(!$res){
                            throw new Exception("加兑换额度失败");
                        }
                        //资金记录
                        $res = $this->db->insert("member_wallet_coin_log",array(
                            'uid'=>$v['uid'],
                            'coin_id'=>$coinInfo['id'],
                            'coin_name'=>$coinInfo['coinname'],
                            'type'=>110,
                            'note'=>'实名认证送认购额度',
                            'value'=> $num1,
                            'balance'=>$wallet['total'],
                            'wallet_id'=>$wallet['id'],
                            'status'=>1,
                            'inputtime'=>time(),
                            'cid'=>$v['uid']
                        ));
                        if(!$res){
                            throw new Exception("加资金记录失败");
                        }
                        //修改状态
                        $res = $this->db->update("member_auth",['upStatus'=>2],['uid'=>$v['uid']]);
                        if(!$res){
                            throw new Exception("修改状态失败");
                        }
                        $content = date("Y-m-d H:i:s").', 成功， 下级实名认证送认购额度：'.json_encode($v).'， 金额：'.$num;
                        file_put_contents(dirname(__FILE__).'/miner_rank.txt', $content.PHP_EOL,FILE_APPEND);
                        return true;
                    }catch (Exception $e){
                        $content = date("Y-m-d H:i:s").', 失败，错误信息：'.$e->getMessage().' 下级实名认证送认购额度：'.json_encode($v).'， 金额：10';
                        file_put_contents(dirname(__FILE__).'/miner_rank.txt', $content.PHP_EOL,FILE_APPEND);
                        return false;
                    }
                });
            }
        }
    }

    //兑换金额解冻
    public function convertedToTotalAction(){
        while (true){
            $status = intval($this->db->get("setting",'content',['type'=>12]));
            if($status != 1){
                sleep(5);
                continue;
            }
            $page = 1;
            while(true){
//            $where['uid'] = 874;
                $where['coin_id'] = $this->coinId;
                $where['converted_amount[>]'] = 0;
                $where['LIMIT'] = [($page - 1)*10,10];
//            $where['OR'] = ['converted_amount[>]'=>0,'convert_amount[>]'=>0];
                $list = $this->db->select("member_wallet",'*',$where);
                if(empty($list)){
                    break;
                }
                foreach($list as $k=>$v){
                    $this->db->action(function() use($v){
                        try{
                            $update = [];
//                       if($v['converted_amount'] > 0){
                            $update['converted_amount[-]'] = $v['converted_amount'];
                            $update['total[+]'] = $v['converted_amount'];
//                       }
//                       if($v['convert_amount'] > 0){
//                           $update['convert_amount[-]'] = $v['convert_amount'];
//                       }
                            $res = $this->db->update("member_wallet",$update,['id'=>$v['id']]);
                            if(!$res){
                                throw new Exception("解冻金额失败");
                            }
                            $content = date("Y-m-d H:i:s").', 解冻已认购金额,成功， '.json_encode($v);
                            file_put_contents(dirname(__FILE__).'/log/unfrozenTransfer.txt', $content.PHP_EOL,FILE_APPEND);
                            return true;
                        }catch (Exception $e){
                            $content = date("Y-m-d H:i:s").', 解冻已认购金额，失败 '.json_encode($v);
                            file_put_contents(dirname(__FILE__).'/log/unfrozenTransferError.txt', $content.PHP_EOL,FILE_APPEND);
                            return false;
                        }
                    });
                }
            }
            //执行完毕 修改状态
            $this->db->update("setting",['content'=>2],['type'=>12]);
            break;
        }
    }

    //已送额度清零
    public function removeAction(){
        while(true){
            $status = intval($this->db->get("setting",'content',['type'=>13]));
            if($status != 1){
                sleep(5);
                continue;
            }
            $this->db->update("member_wallet",['convert_amount'=>0],['coin_id'=>6,'convert_amount[>]'=>0]);
            //执行完毕 修改状态
            $this->db->update("setting",['content'=>2],['type'=>13]);
            break;
        }
    }

    //2020-05-19  今天兑换的退回去
    public function returnCoinAction(){
        $start = strtotime(date("Y-m-d"));
        $where['coin_id'] = 5;
        $where['type'] = 112;
        $where['inputtime[>]'] = $start;
        $list = $this->db->select("member_wallet_coin_log",'*',$where);
        foreach ($list as $v){
            //是否已经退回
            $log = $this->db->get("member_wallet_coin_log",'*',['type'=>131,'uid'=>$v['uid']]);
            if($log){
                continue;
            }
            $wallet = $this->db->get("member_wallet",'*',['uid'=>$v['uid'],'coin_id'=>$v['coin_id']]);
            if(!$wallet){
                continue;
            }
            $this->db->action(function() use($v,$wallet){
               try{
                   //减金额
                   $res = $this->db->update("member_wallet",['total[+]'=>abs($v['value'])],['id'=>$wallet['id']]);
                   if(!$res){
                       throw new Exception("退回金额失败");
                   }
                   //记录
                   $res = $this->db->insert("member_wallet_coin_log",array(
                       'uid'=>$v['uid'],
                       'coin_id'=>$v['coin_id'],
                       'coin_name'=>$v['coin_name'],
                       'wallet_id'=>$wallet['id'],
                       'value' => abs($v['value']),
                       'balance'=>$wallet['total'],
                       'cid'=>$v['id'],
                       'inputtime'=>time(),
                       'type'=>131,
                       'note'=>'认购退回'
                   ));
                   if(!$res){
                       throw new Exception("记录失败");
                   }
                   //返回额度
                   $wa1 = $this->db->get("member_wallet",'*',['uid'=>$v['uid'],'coin_id'=>6]);
                   $res = $this->db->update("member_wallet",['convert_amount[+]'=>$wa1['converted_amount'],'converted_amount[-]'=>$wa1['converted_amount']],['id'=>$wa1['id']]);
                   if(!$res){
                       throw new Exception("加回额度失败");
                   }
                   return true;
               }catch (Exception $e){
                   return false;
               }
            });
            exit;
        }
        print_r($list);exit;
    }

}
