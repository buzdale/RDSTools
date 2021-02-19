#### describeRDS.go


### Uses AWS API to get info about database instances.
export AWS_REGION= whatever your region is

USAGE: 

        describeRDS -c	
        Calculates average CPU usage last 5 minutes

        describeRDS -d <dbInstance>
        Describes Database Instance - If blank describes all in the Region.
        
       describeRDS -l 
        Lists databaseInstances in region
        
       describeRDS -f
        Shows free storage in GBytes for all non-aurora dbInstances.
         
