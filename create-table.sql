-- connect to database: mysql -u root -p
-- create database
-- CREATE DATABASE cms;
-- create database user
-- GRANT ALL ON cms.* TO cms@localhost IDENTIFIED BY 'cms123';

-- create `pages` table
CREATE TABLE `pages` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `page_guid` varchar(256) NOT NULL DEFAULT '',
  `page_title` varchar(256) DEFAULT NULL,
  `page_content` mediumtext,
  `page_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `page_guid` (`page_guid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

-- insert 1 row
INSERT INTO `pages` (`id`, `page_guid`, `page_title`, `page_content`, `page_date`) VALUES (NULL, 'hello-world', 'Hello, World', 'I\'m so glad you found this page!  It\'s been sitting patiently on the Internet for some time, just waiting for a visitor.', CURRENT_TIMESTAMP);

-- insert second row
INSERT INTO `pages` (`id`, `page_guid`, `page_title`, `page_content`, `page_date`) VALUES (2, 'a-new-blog', 'A New Blog', 'I hope you enjoyed the last blog!  Well brace yourself, because my latest blog is even <i>better</i> than the last!', '2015-04-29 02:16:19');

-- insert another row
INSERT INTO `pages` (`id`, `page_guid`, `page_title`, `page_content`, `page_date`) VALUES (3, 'lorem-ipsum', 'Lorem Ipsum', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam vel dapibus enim. Nulla a accumsan nisi. Cras vel nunc ullamcorper, fringilla justo vel, molestie magna. Suspendisse rhoncus vulputate tortor ac pulvinar. Mauris laoreet viverra semper. Mauris faucibus non nisl ut semper. Vestibulum in tellus sed ligula lacinia pulvinar eget ac nisi. Phasellus vitae ornare velit. Pellentesque et ipsum nibh. Vivamus elementum egestas justo. Nulla non diam et nibh suscipit ultricies. Ut suscipit risus nec libero cursus, at molestie mauris tincidunt. Etiam et orci justo.', '2015-05-06 04:09:45');

-- create `comments` table
CREATE TABLE `comments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `page_id` int(11) NOT NULL,
  `comment_guid` varchar(256) DEFAULT NULL,
  `comment_name` varchar(64) DEFAULT NULL,
  `comment_email` varchar(128) DEFAULT NULL,
  `comment_text` mediumtext,
  `comment_date` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `page_id` (`page_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;