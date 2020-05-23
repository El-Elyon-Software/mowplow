SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mowplow
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mowplow
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mowplow`;
USE `mowplow` ;

-- -----------------------------------------------------
-- Table `mowplow`.`city_state_country`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`city_state_country` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`city_state_country` (
  `postal_code` VARCHAR(15) NOT NULL,
  `city` VARCHAR(100) NOT NULL,
  `state` VARCHAR(100) NOT NULL,
  `state_code` VARCHAR(2) NOT NULL,
  `country` VARCHAR(85) NOT NULL,
  `country_code` VARCHAR(45) NOT NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`postal_code`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mowplow`.`end_customer`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`end_customer` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`end_customer` (
  `end_customer_id` INT NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `business_name` VARCHAR(45) NULL,
  `address_1` VARCHAR(100) NOT NULL,
  `address_2` VARCHAR(45) NULL,
  `postal_code` VARCHAR(15) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `mobile` VARCHAR(15) NOT NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`end_customer_id`),
  CONSTRAINT `fk_end_customer_postal_code`
    FOREIGN KEY (`postal_code`)
    REFERENCES `mowplow`.`city_state_country` (`postal_code`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `mowplow`.`service_provider`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`service_provider` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`service_provider` (
  `service_provider_id` INT NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `business_name` VARCHAR(45) NULL,
  `address_1` VARCHAR(100) NOT NULL,
  `address_2` VARCHAR(45) NULL,
  `postal_code` VARCHAR(15) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `mobile` VARCHAR(15) NOT NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`service_provider_id`),
  CONSTRAINT `fk_service_provider_postal_code`
    FOREIGN KEY (`postal_code`)
    REFERENCES `mowplow`.`city_state_country` (`postal_code`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `mowplow`.`service`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`service` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`service` (
  `service_id` INT NOT NULL AUTO_INCREMENT,
  `service_name` VARCHAR(45) NOT NULL,
  `description` VARCHAR(255) NOT NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`service_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mowplow`.`service_provider_service`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`service_provider_service` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`service_provider_service` (
  `service_provider_service_id` INT NOT NULL AUTO_INCREMENT,
  `service_provider_id` INT NULL,
  `service_id` INT NULL,
  `description` VARCHAR(255) NOT NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`service_provider_service_id`),
  CONSTRAINT `fk_service_provider_service_service_provider_id`
    FOREIGN KEY (`service_provider_id`)
    REFERENCES `mowplow`.`service_provider` (`service_provider_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_service_provider_service_service_id`
    FOREIGN KEY (`service_id`)
    REFERENCES `mowplow`.`service` (`service_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table `mowplow`.`end_customer_service`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`end_customer_service` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`end_customer_service` (
  `end_customer_service_id` INT NOT NULL AUTO_INCREMENT,
  `end_customer_id` INT NULL,
  `service_provider_service_id` INT NULL,
  `description` VARCHAR(255) NULL,
  `estimated_job_length` DOUBLE(8,2) NOT NULL,
  `contract_start_date` DATETIME NOT NULL,
  `contract_end_date` DATETIME NOT NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`end_customer_service_id`),
  CONSTRAINT `fk_end_customer_service_end_customer_id`
    FOREIGN KEY (`end_customer_id`)
    REFERENCES `mowplow`.`end_customer` (`end_customer_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_end_customer_service_provider_service_id`
    FOREIGN KEY (`service_provider_service_id`)
    REFERENCES `mowplow`.`service_provider_service` (`service_provider_service_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mowplow`.`service_schedule`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `mowplow`.`service_schedule` ;

CREATE TABLE IF NOT EXISTS `mowplow`.`service_schedule` (
  `end_customer_service_id` INT NOT NULL,
  `date` DATETIME NOT NULL,
  `place_in_line` INT NOT NULL,
  `actual_job_length` DECIMAL(8,2) NULL,
  `estimated_start_time` TIME NULL,
  `actual_start_time` TIME NULL,
  `completed_time` TIME NULL,
  `date_added` DATETIME NOT NULL DEFAULT NOW(),
  `date_modified` DATETIME NOT NULL DEFAULT NOW(),
  `date_deleted` DATETIME NULL,
  PRIMARY KEY (`end_customer_service_id`, `date`),
  CONSTRAINT `fk_service_schedule_end_customer_service_id`
    FOREIGN KEY (`end_customer_service_id`)
    REFERENCES `mowplow`.`end_customer_service` (`end_customer_service_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
