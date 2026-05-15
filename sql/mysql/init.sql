-- WVP-PRO-GO MySQL Schema
-- Auto-generated from Java WVP-PRO v2.7.4
-- Note: GORM AutoMigrate handles table creation at runtime

-- Storage for GB28181 device basic info and online status
CREATE TABLE IF NOT EXISTS wvp_device (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL UNIQUE COMMENT 'GB device ID',
    name VARCHAR(255) COMMENT 'Device name',
    manufacturer VARCHAR(255) COMMENT 'Manufacturer',
    model VARCHAR(255) COMMENT 'Model',
    firmware VARCHAR(255) COMMENT 'Firmware version',
    transport VARCHAR(50) COMMENT 'SIP transport (TCP/UDP)',
    stream_mode VARCHAR(50) COMMENT 'Stream pull mode',
    on_line BOOLEAN DEFAULT FALSE COMMENT 'Online status',
    ip VARCHAR(50) COMMENT 'Device IP',
    port INT COMMENT 'SIP port',
    expires INT COMMENT 'Registration expiry',
    host_address VARCHAR(50) COMMENT 'Host address',
    charset VARCHAR(50) COMMENT 'Charset',
    ssrc_check BOOLEAN DEFAULT FALSE COMMENT 'SSRC check',
    geo_coord_sys VARCHAR(50) COMMENT 'Geo coordinate system',
    media_server_id VARCHAR(50) DEFAULT 'auto' COMMENT 'Media server ID',
    custom_name VARCHAR(255) COMMENT 'Custom display name',
    sdp_ip VARCHAR(50) COMMENT 'SDP IP',
    local_ip VARCHAR(50) COMMENT 'Local IP',
    password VARCHAR(255) COMMENT 'Auth password',
    as_message_channel BOOLEAN DEFAULT FALSE COMMENT 'As message channel',
    heart_beat_interval INT COMMENT 'Heartbeat interval',
    heart_beat_count INT COMMENT 'Heartbeat failure count',
    position_capability INT COMMENT 'Position capability',
    channel_count INT COMMENT 'Channel count',
    subscribe_cycle_for_catalog INT DEFAULT 0 COMMENT 'Catalog subscribe cycle',
    subscribe_cycle_for_mobile_position INT DEFAULT 0 COMMENT 'Mobile position subscribe cycle',
    mobile_position_submission_interval INT DEFAULT 5 COMMENT 'Mobile position report interval',
    subscribe_cycle_for_alarm INT DEFAULT 0 COMMENT 'Alarm subscribe cycle',
    broadcast_push_after_ack BOOLEAN DEFAULT FALSE COMMENT 'Push stream after ACK',
    server_id VARCHAR(50) COMMENT 'Server ID',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='GB28181 device';

-- Device channel information
CREATE TABLE IF NOT EXISTS wvp_device_channel (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    gb_id BIGINT AUTO_INCREMENT COMMENT 'GB ID',
    device_id VARCHAR(50) COMMENT 'Device ID',
    name VARCHAR(255) COMMENT 'Channel name',
    gb_device_id VARCHAR(50) COMMENT 'GB device ID',
    gb_name VARCHAR(255) COMMENT 'GB name',
    gb_manufacturer VARCHAR(50) COMMENT 'GB manufacturer',
    gb_model VARCHAR(50) COMMENT 'GB model',
    gb_owner VARCHAR(50) COMMENT 'GB owner',
    gb_civil_code VARCHAR(50) COMMENT 'GB civil code',
    gb_block VARCHAR(50) COMMENT 'GB block',
    gb_address VARCHAR(50) COMMENT 'GB address',
    gb_parental INT COMMENT 'GB parental',
    gb_parent_id VARCHAR(50) COMMENT 'GB parent ID',
    gb_safety_way INT COMMENT 'GB safety way',
    gb_register_way INT COMMENT 'GB register way',
    gb_cert_num VARCHAR(50) COMMENT 'GB cert number',
    gb_certifiable INT COMMENT 'GB certifiable',
    gb_err_code INT COMMENT 'GB error code',
    gb_end_time VARCHAR(50) COMMENT 'GB end time',
    gb_secrecy INT COMMENT 'GB secrecy',
    gb_ip_address VARCHAR(50) COMMENT 'GB IP address',
    gb_port INT COMMENT 'GB port',
    gb_password VARCHAR(255) COMMENT 'GB password',
    gb_status VARCHAR(50) COMMENT 'GB status',
    gb_longitude DOUBLE COMMENT 'GB longitude',
    gb_latitude DOUBLE COMMENT 'GB latitude',
    gps_altitude DOUBLE COMMENT 'GPS altitude',
    gps_speed DOUBLE COMMENT 'GPS speed',
    gps_direction DOUBLE COMMENT 'GPS direction',
    gps_time VARCHAR(50) COMMENT 'GPS time',
    gb_business_group_id VARCHAR(50) COMMENT 'GB business group ID',
    gb_ptz_type INT COMMENT 'GB PTZ type',
    gb_position_type INT COMMENT 'GB position type',
    gb_room_type INT COMMENT 'GB room type',
    gb_use_type INT COMMENT 'GB use type',
    gb_supply_light_type INT COMMENT 'GB supply light type',
    gb_direction_type INT COMMENT 'GB direction type',
    gb_resolution VARCHAR(50) COMMENT 'GB resolution',
    gb_download_speed VARCHAR(50) COMMENT 'GB download speed',
    gb_svc_space_support_mod INT COMMENT 'GB SVC space support',
    gb_svc_time_support_mode INT COMMENT 'GB SVC time support',
    record_plan VARCHAR(50) COMMENT 'Record plan',
    data_type INT DEFAULT 1 COMMENT 'Data type (1=GB28181)',
    data_device_id INT COMMENT 'Data device ID',
    stream_identification VARCHAR(50) COMMENT 'Stream identification',
    enable_broadcast INT DEFAULT 0 COMMENT 'Enable broadcast',
    map_level INT DEFAULT 0 COMMENT 'Map level',
    parental INT COMMENT 'Parental',
    parent_id VARCHAR(50) COMMENT 'Parent ID',
    safety_way INT COMMENT 'Safety way',
    register_way INT COMMENT 'Register way',
    cert_num VARCHAR(50) COMMENT 'Cert number',
    certifiable INT COMMENT 'Certifiable',
    err_code INT COMMENT 'Error code',
    end_time VARCHAR(50) COMMENT 'End time',
    secrecy INT COMMENT 'Secrecy',
    ip_address VARCHAR(50) COMMENT 'IP address',
    port INT COMMENT 'Port',
    password VARCHAR(255) COMMENT 'Password',
    status VARCHAR(50) COMMENT 'Status',
    longitude DOUBLE COMMENT 'Longitude',
    latitude DOUBLE COMMENT 'Latitude',
    ptz_type INT COMMENT 'PTZ type',
    position_type INT COMMENT 'Position type',
    room_type INT COMMENT 'Room type',
    use_type INT COMMENT 'Use type',
    supply_light_type INT COMMENT 'Supply light type',
    direction_type INT COMMENT 'Direction type',
    resolution VARCHAR(50) COMMENT 'Resolution',
    manufacturer VARCHAR(50) COMMENT 'Manufacturer',
    model VARCHAR(50) COMMENT 'Model',
    owner VARCHAR(50) COMMENT 'Owner',
    civil_code VARCHAR(50) COMMENT 'Civil code',
    block VARCHAR(50) COMMENT 'Block',
    address VARCHAR(50) COMMENT 'Address',
    stream_id VARCHAR(50) COMMENT 'Stream ID',
    INDEX idx_device_id (device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Device channel';

-- Cascade platform info
CREATE TABLE IF NOT EXISTS wvp_platform (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    enable INT DEFAULT 0 COMMENT 'Enable status',
    name VARCHAR(255) COMMENT 'Platform name',
    server_gb_id VARCHAR(50) UNIQUE COMMENT 'Server GB ID',
    server_gb_domain VARCHAR(50) COMMENT 'Server GB domain',
    server_ip VARCHAR(50) COMMENT 'Server IP',
    server_port INT COMMENT 'Server port',
    device_gb_id VARCHAR(50) COMMENT 'Device GB ID',
    device_ip VARCHAR(50) COMMENT 'Device IP',
    device_port INT COMMENT 'Device port',
    username VARCHAR(255) COMMENT 'Username',
    password VARCHAR(255) COMMENT 'Password',
    expires INT COMMENT 'Expires',
    keep_timeout INT COMMENT 'Keep timeout',
    transport VARCHAR(50) COMMENT 'Transport',
    character_set VARCHAR(50) COMMENT 'Character set',
    ptz INT COMMENT 'PTZ',
    rtcp INT COMMENT 'RTCP',
    status VARCHAR(50) COMMENT 'Status',
    channel_count INT COMMENT 'Channel count',
    catalog_subscribe INT COMMENT 'Catalog subscribe',
    alarm_subscribe INT COMMENT 'Alarm subscribe',
    mobile_position_subscribe INT COMMENT 'Mobile position subscribe',
    catalog_group INT COMMENT 'Catalog group',
    as_message_channel BOOLEAN DEFAULT FALSE COMMENT 'As message channel',
    send_stream_ip VARCHAR(50) COMMENT 'Send stream IP',
    auto_push_channel BOOLEAN DEFAULT FALSE COMMENT 'Auto push channel',
    catalog_with_platform INT COMMENT 'Catalog with platform',
    catalog_with_group INT COMMENT 'Catalog with group',
    catalog_with_region INT COMMENT 'Catalog with region',
    civil_code VARCHAR(50) COMMENT 'Civil code',
    manufacturer VARCHAR(255) COMMENT 'Manufacturer',
    model VARCHAR(255) COMMENT 'Model',
    address VARCHAR(255) COMMENT 'Address',
    register_way INT COMMENT 'Register way',
    secrecy INT COMMENT 'Secrecy',
    server_id VARCHAR(50) COMMENT 'Server ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Cascade platform';

-- Platform-channel relationship
CREATE TABLE IF NOT EXISTS wvp_platform_channel (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform_id BIGINT COMMENT 'Platform ID',
    channel_id BIGINT COMMENT 'Channel ID',
    INDEX idx_platform_id (platform_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Platform channel mapping';

-- Business group
CREATE TABLE IF NOT EXISTS wvp_common_group (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(50) UNIQUE COMMENT 'Device ID',
    name VARCHAR(255) COMMENT 'Group name',
    parent_id BIGINT COMMENT 'Parent ID',
    parent_device_id VARCHAR(50) COMMENT 'Parent device ID',
    business_group VARCHAR(50) COMMENT 'Business group',
    civil_code VARCHAR(50) COMMENT 'Civil code',
    alias VARCHAR(255) COMMENT 'Alias',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Business group';

-- Administrative region
CREATE TABLE IF NOT EXISTS wvp_common_region (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(50) UNIQUE COMMENT 'Device ID',
    name VARCHAR(255) COMMENT 'Region name',
    parent_id BIGINT COMMENT 'Parent ID',
    parent_device_id VARCHAR(50) COMMENT 'Parent device ID',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Administrative region';

-- Stream proxy
CREATE TABLE IF NOT EXISTS wvp_stream_proxy (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type INT COMMENT 'Type (0=pull proxy)',
    app VARCHAR(255) COMMENT 'App name',
    stream VARCHAR(255) COMMENT 'Stream name',
    url VARCHAR(255) COMMENT 'Stream URL',
    ffmpeg_cmd TEXT COMMENT 'FFmpeg command',
    enable_audio BOOLEAN COMMENT 'Enable audio',
    enable_mp4 BOOLEAN COMMENT 'Enable MP4 recording',
    enable BOOLEAN COMMENT 'Enable status',
    timeout INT COMMENT 'Timeout (seconds)',
    pulling BOOLEAN COMMENT 'Is pulling',
    enable_remove_key BOOLEAN COMMENT 'Enable remove key',
    remove_key VARCHAR(255) COMMENT 'Remove key',
    media_server_id VARCHAR(50) COMMENT 'Media server ID',
    channel_id BIGINT COMMENT 'Channel ID',
    device_id VARCHAR(50) COMMENT 'Device ID',
    name VARCHAR(255) COMMENT 'Name',
    description VARCHAR(255) COMMENT 'Description',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Stream proxy';

-- Stream push
CREATE TABLE IF NOT EXISTS wvp_stream_push (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    app VARCHAR(255) COMMENT 'App name',
    stream VARCHAR(255) COMMENT 'Stream name',
    total_origin INT COMMENT 'Total origin',
    enable_audio BOOLEAN COMMENT 'Enable audio',
    enable_mp4 BOOLEAN COMMENT 'Enable MP4 recording',
    pushing BOOLEAN COMMENT 'Is pushing',
    media_server_id VARCHAR(50) COMMENT 'Media server ID',
    channel_id BIGINT COMMENT 'Channel ID',
    device_id VARCHAR(50) COMMENT 'Device ID',
    name VARCHAR(255) COMMENT 'Name',
    description VARCHAR(255) COMMENT 'Description',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Stream push';

-- Media server
CREATE TABLE IF NOT EXISTS wvp_media_server (
    id VARCHAR(50) PRIMARY KEY,
    ip VARCHAR(50) COMMENT 'IP',
    hook_ip VARCHAR(50) COMMENT 'Hook IP',
    sdp_ip VARCHAR(50) COMMENT 'SDP IP',
    stream_ip VARCHAR(50) COMMENT 'Stream IP',
    http_port INT COMMENT 'HTTP port',
    http_ssl_port INT COMMENT 'HTTP SSL port',
    rtmp_port INT COMMENT 'RTMP port',
    rtmp_ssl_port INT COMMENT 'RTMP SSL port',
    flv_port INT COMMENT 'FLV port',
    flv_ssl_port INT COMMENT 'FLV SSL port',
    mp4_port INT COMMENT 'MP4 port',
    ws_flv_port INT COMMENT 'WS FLV port',
    ws_flv_ssl_port INT COMMENT 'WS FLV SSL port',
    rtsp_port INT COMMENT 'RTSP port',
    rtsp_ssl_port INT COMMENT 'RTSP SSL port',
    rtp_proxy_port INT COMMENT 'RTP proxy port',
    jtt_proxy_port INT COMMENT 'JTT proxy port',
    auto_config BOOLEAN COMMENT 'Auto config',
    secret VARCHAR(255) COMMENT 'Secret',
    hook_alive_interval DOUBLE COMMENT 'Hook alive interval',
    rtp_enable BOOLEAN COMMENT 'RTP enable',
    status BOOLEAN COMMENT 'Status',
    rtp_port_range VARCHAR(50) COMMENT 'RTP port range',
    send_rtp_port_range VARCHAR(50) COMMENT 'Send RTP port range',
    record_assist_port INT COMMENT 'Record assist port',
    default_server BOOLEAN COMMENT 'Default server',
    record_day INT COMMENT 'Record days',
    record_path VARCHAR(255) COMMENT 'Record path',
    type VARCHAR(50) COMMENT 'Type (zlm/abl)',
    transcode_suffix VARCHAR(50) COMMENT 'Transcode suffix',
    server_id VARCHAR(50) COMMENT 'Server ID',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time',
    last_keepalive_time VARCHAR(50) COMMENT 'Last keepalive time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Media server';

-- Mobile position
CREATE TABLE IF NOT EXISTS wvp_device_mobile_position (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL COMMENT 'Device ID',
    channel_id VARCHAR(50) NOT NULL COMMENT 'Channel ID',
    device_name VARCHAR(255) COMMENT 'Device name',
    time VARCHAR(50) COMMENT 'Report time',
    longitude DOUBLE COMMENT 'Longitude',
    latitude DOUBLE COMMENT 'Latitude',
    altitude DOUBLE COMMENT 'Altitude',
    speed DOUBLE COMMENT 'Speed',
    direction DOUBLE COMMENT 'Direction',
    report_source VARCHAR(50) COMMENT 'Report source',
    create_time VARCHAR(50) COMMENT 'Create time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Mobile position';

-- Device alarm
CREATE TABLE IF NOT EXISTS wvp_device_alarm (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL COMMENT 'Device ID',
    channel_id VARCHAR(50) NOT NULL COMMENT 'Channel ID',
    alarm_priority VARCHAR(50) COMMENT 'Alarm priority',
    alarm_method VARCHAR(50) COMMENT 'Alarm method',
    alarm_time VARCHAR(50) COMMENT 'Alarm time',
    alarm_description VARCHAR(255) COMMENT 'Alarm description',
    longitude DOUBLE COMMENT 'Longitude',
    latitude DOUBLE COMMENT 'Latitude',
    alarm_type VARCHAR(50) COMMENT 'Alarm type',
    create_time VARCHAR(50) NOT NULL COMMENT 'Create time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Device alarm';

-- JT1078 terminal
CREATE TABLE IF NOT EXISTS wvp_jt_terminal (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL UNIQUE COMMENT 'Phone number',
    plate_number VARCHAR(20) COMMENT 'Plate number',
    plate_color INT COMMENT 'Plate color',
    sim_card_id VARCHAR(20) COMMENT 'SIM card ID',
    terminal_id VARCHAR(20) COMMENT 'Terminal ID',
    terminal_model VARCHAR(20) COMMENT 'Terminal model',
    manufacturer_id VARCHAR(20) COMMENT 'Manufacturer ID',
    province_id INT COMMENT 'Province ID',
    city_id INT COMMENT 'City ID',
    online BOOLEAN DEFAULT FALSE COMMENT 'Online status',
    media_server_id VARCHAR(50) COMMENT 'Media server ID',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='JT1078 terminal';

-- JT1078 channel
CREATE TABLE IF NOT EXISTS wvp_jt_channel (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL COMMENT 'Phone number',
    channel_id VARCHAR(20) NOT NULL COMMENT 'Channel ID',
    channel_name VARCHAR(255) COMMENT 'Channel name',
    channel_type INT COMMENT 'Channel type'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='JT1078 channel';

-- User
CREATE TABLE IF NOT EXISTS wvp_user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE COMMENT 'Username',
    password VARCHAR(255) NOT NULL COMMENT 'Password',
    name VARCHAR(255) COMMENT 'Name',
    phone VARCHAR(20) COMMENT 'Phone',
    email VARCHAR(255) COMMENT 'Email',
    enable BOOLEAN DEFAULT TRUE COMMENT 'Enable status',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User';

-- Role
CREATE TABLE IF NOT EXISTS wvp_user_role (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE COMMENT 'Role name',
    description VARCHAR(255) COMMENT 'Description',
    create_time VARCHAR(50) COMMENT 'Create time',
    update_time VARCHAR(50) COMMENT 'Update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User role';
