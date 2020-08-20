#### describeRDS.go


### Uses AWS API to get info about database instances.
export AWS_REGION= whatever your region is

USAGE: 

        describeRDS -d <dbInstance>
        describes Database Instance - If blank describes all in the Region.
        
       describeRDS -l 
        lists databaseInstances in region
        
       describeRDS -f
        Shows free storage in GBytes for all non-aurora dbInstances.
         
