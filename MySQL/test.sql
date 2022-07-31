Use syncer;

-- syncer.Partner definition

CREATE TABLE IF NOT EXISTS `Partner` (
  `ID` varchar(129) NOT NULL,
  `Name` varchar(100) NOT NULL,
  `URL` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `Partner_ID_IDX` (`ID`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- syncer.Sync definition

CREATE TABLE IF NOT EXISTS `Sync` (
  `PartnerID` varchar(129) NOT NULL,
  `PartnerUserID` varchar(129) DEFAULT NULL,
  `OtherPartnerID` varchar(129) DEFAULT NULL,
  `OtherPartnerUserID` varchar(129) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
